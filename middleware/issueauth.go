package middleware

import (
	"encoding/base64"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gitlab.orayer.com/golang/issue/library/container"
	"net/http"
)

type AuthRecord struct {
	Action string `json:"action" binding:"required"`
	UniqId string `json:"uniqid" binding:"required"`
	Secret string `json:"secret" binding:"required"`
}

type StorageRecord struct {
	CryptoType string `json:"cryptotype"`
	Action     string `json:"action"`
	UniqId     string `json:"uniqid"`
	Secret     string `json:"secret"`
}

const (
	CryptoTypeClear  = "CLEAR"
	CryptoTypeTicket = "TICKET"

	AuthInfoKey = "_AUTH_"
)

func IssueAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.Query("key")
		if key == "" {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		baseData, err := base64.StdEncoding.DecodeString(key)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		var authRecord AuthRecord
		var storageRecord StorageRecord

		if err := json.Unmarshal(baseData, &authRecord); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		authKey := "issue_auth_key_" + authRecord.Action + "_" + authRecord.UniqId
		storageData, err := container.Mgr.Storager.Get(authKey)
		if err != nil {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		if err := json.Unmarshal(storageData, &storageRecord); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		switch storageRecord.CryptoType {
		case CryptoTypeClear:
		case CryptoTypeTicket:
			fallthrough
		default:
			if storageRecord.Secret != authRecord.Secret {
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}

			_ = container.Mgr.Storager.Delete(authKey)
		}

		c.Set(AuthInfoKey, storageRecord)
		c.Next()
	}
}
