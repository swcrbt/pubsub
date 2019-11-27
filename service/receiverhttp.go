package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.orayer.com/golang/errors"
	"gitlab.orayer.com/golang/issue/library/container"
	"gitlab.orayer.com/golang/server"
	"net/http"
)

type ReleaseBody struct {
	// 推送ID
	Topics []string `from:"topics" json:"topics" binding:"required"`

	// 内容行为
	Action string `from:"action" json:"action" binding:"required"`

	// 内容
	Body interface{} `from:"body" json:"body"`
}

type HttpReceiver struct {
	handler func(c *gin.Context)
	server *server.HttpServer
}

func NewHttpReceiver() *HttpReceiver {
	return &HttpReceiver{
		handler: receiverHandler,
	}
}

func (rec *HttpReceiver) Run() error {
	rev := server.NewHttpServer()

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

	rec.server = rev

	return nil
}

func (rec *HttpReceiver) GetName() string {
	return "receiver-http"
}

func (rec *HttpReceiver) Stop() error  {
	if rec.server != nil {
		return rec.server.Shutdown()
	}
	return nil
}

func receiverHandler(c *gin.Context) {
	var params ReleaseBody

	defer c.Request.Body.Close()

	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, errors.New("receiverhttp/params_error", "invalid request payload"))
		return
	}

	result := container.Mgr.Dispatcher.Publish(params.Topics, params.Action, params.Body)

	fmt.Println(result)

	c.JSON(http.StatusOK, result)
	return
}

