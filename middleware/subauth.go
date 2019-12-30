package middleware

import (
	"encoding/base64"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gitlab.orayer.com/golang/pubsub/library/container"
	"net/http"
)

type AuthRecord struct {
	Key    string `json:"key" binding:"required"`
	Secret string `json:"secret" binding:"required"`
}

type StorageRecord struct {
	// 加密类型，CLEAR：存在即可，TICKET：一次性登录凭证
	CryptoType string `json:"cryptotype"`

	// Id，存在时会出现互踢
	ID string `json:"id"`

	// 默认订阅主题
	Topics []string `json:"topics"`

	// 密钥
	Secret string `json:"secret"`
}

const (
	CryptoTypeClear  = "CLEAR"
	CryptoTypeTicket = "TICKET"

	AuthKeyPrefix = "sub_auth_key_"
	AuthInfoKey = "_AUTH_"
)

func SubAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.Query("key")
		if key == "" {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		baseData, err := base64.StdEncoding.DecodeString(key)
		if err != nil {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		var authRecord AuthRecord
		var storageRecord StorageRecord

		if err := json.Unmarshal(baseData, &authRecord); err != nil {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		authKey := AuthKeyPrefix + authRecord.Key
		storageData, err := container.Mgr.Storager.Get(authKey)
		if err != nil {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		if err := json.Unmarshal(storageData, &storageRecord); err != nil {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		switch storageRecord.CryptoType {
		case CryptoTypeClear:
		case CryptoTypeTicket:
			fallthrough
		default:
			if storageRecord.Secret != authRecord.Secret {
				c.AbortWithStatus(http.StatusForbidden)
				return
			}

			_ = container.Mgr.Storager.Delete(authKey)
		}

		c.Set(AuthInfoKey, storageRecord)
		c.Next()
	}
}
