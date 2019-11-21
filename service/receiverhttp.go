package service

import (
	"github.com/gin-gonic/gin"
	"gitlab.orayer.com/golang/errors"
	"gitlab.orayer.com/golang/server"
	"go-issued-service/library/container"
	"net/http"
)

type ReleaseBody struct {
	// 推送类型
	Action string `from:"action" json:"action" binding:"required"`

	// 推送ID
	UniqIds []string `from:"uniqids" json:"uniqids" binding:"required"`

	// 推送内容
	Data interface{} `from:"data" json:"data" binding:"required"`
}

type HttpReceiver struct {
	handler func(c *gin.Context)
}

func NewHttpReceiver() *HttpReceiver {
	return &HttpReceiver{
		handler: receiverHandler,
	}
}

func (rec *HttpReceiver) Run() error {
	rev, err := server.Endless().NewHttpServer()
	if (err != nil) {
		return err
	}

	gin.SetMode(container.Mgr.Config.Server.Mode)
	rev.Router.Use(gin.Logger(), gin.Recovery())

	rev.Router.POST("/release", rec.handler)
	rev.Port = container.Mgr.Config.Server.ReceiverHttp.Port

	go func() {
		container.Mgr.Logger.Printf("\"%s\" Server Run At: \"%d\"\n", rec.GetName(), container.Mgr.Config.Server.ReceiverHttp.Port)

		if err := rev.Start(); err != nil {
			container.Mgr.Logger.Printf("\"%s\" Server error: %v\n", rec.GetName(), err)
		}
	}()

	return nil
}

func (rec *HttpReceiver) GetName() string {
	return "receiver-http"
}

func receiverHandler(c *gin.Context) {
	var params ReleaseBody

	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, errors.New("receiverhttp/params_error", "invalid request payload"))
		return
	}

	container.Mgr.Dispatcher.Push(params.Action, params.UniqIds, params.Data)
}
