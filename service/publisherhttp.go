package service

import (
	"github.com/gin-gonic/gin"
	"gitlab.orayer.com/golang/errors"
	"gitlab.orayer.com/golang/pubsub/library/container"
	"gitlab.orayer.com/golang/server"
	"net/http"
)

type PublishBody struct {
	// 推送ID
	Topics []string `from:"topics" json:"topics" binding:"required"`

	// 内容行为
	Action string `from:"action" json:"action" binding:"required"`

	// 内容
	Body interface{} `from:"body" json:"body"`
}

type HttpPublisher struct {
	handler func(c *gin.Context)
	server *server.HttpServer
}

func NewHttpPublisher() *HttpPublisher {
	return &HttpPublisher{
		handler: publisherHandler,
	}
}

func (rec *HttpPublisher) Run() error {
	rev := server.NewHttpServer()

	rev.Router.Use(gin.Logger(), gin.Recovery())

	rev.Router.POST("/publish", rec.handler)
	rev.Port = container.Mgr.Config.Server.PublisherHttp.Port

	go func() {
		container.Mgr.Logger.Printf("\"%s\" Server Run At: \"%d\"\n", rec.GetName(), container.Mgr.Config.Server.PublisherHttp.Port)

		if err := rev.Start(); err != nil {
			container.Mgr.Logger.Printf("\"%s\" Server error: %v\n", rec.GetName(), err)
		}
	}()

	rec.server = rev

	return nil
}

func (rec *HttpPublisher) GetName() string {
	return "publisher-http"
}

func (rec *HttpPublisher) Stop() error  {
	if rec.server != nil {
		return rec.server.Shutdown()
	}
	return nil
}

func publisherHandler(c *gin.Context) {
	var params PublishBody

	defer c.Request.Body.Close()

	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, errors.New("Publisherhttp/params_error", "invalid request payload"))
		return
	}

	/*resp := map[string]map[string]bool{}
	result := container.Mgr.Dispatcher.Publish(params.Topics, params.Action, params.Body)

	result.Range(func(topic, res interface{}) bool {

		if _, ok := resp[topic.(string)]; !ok {
			resp[topic.(string)] = map[string]bool{}
		}

		res.(*sync.Map).Range(func(cid, isReply interface{}) bool {
			resp[topic.(string)][cid.(string)] = isReply.(bool)
			return true
		})

		return true
	})

	c.JSON(http.StatusOK, resp)*/
	container.Mgr.Dispatcher.Publish(params.Topics, params.Action, params.Body)
	c.AbortWithStatus(http.StatusNoContent)
	return
}
