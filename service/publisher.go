package service

import (
	"context"
	"github.com/gin-gonic/gin"
	"gitlab.orayer.com/golang/errors"
	"gitlab.orayer.com/golang/pubsub/library/container"
	"net/http"
	"time"
)

type PublishMessage struct {
	// 推送ID
	Topics []string `from:"topics" json:"topics" binding:"required"`

	// 内容行为
	Action string `from:"action" json:"action" binding:"required"`

	// 内容
	Body interface{} `from:"body" json:"body"`
}

type Publisher struct {
	handler  func(c *gin.Context)
	server   *http.Server
}

func NewPublisher() *Publisher {
	return &Publisher{
		handler:  publisherHandler,
	}
}

func (pub *Publisher) Run() error {
	router := gin.Default()
	router.POST("/publish", pub.handler)

	rev := &http.Server{
		Addr:    container.Mgr.Config.Server.Publisher.Address,
		Handler: router,
	}

	go func() {
		container.Mgr.Logger.Printf("\"%s\" server run at: \"%s\"\n", pub.GetName(), rev.Addr)

		if err := rev.ListenAndServe(); err != nil {
			container.Mgr.Logger.Printf("\"%s\" server stop at \"%s\": %v\n", pub.GetName(), rev.Addr, err)
		}
	}()

	pub.server = rev

	return nil
}

func (pub *Publisher) GetName() string {
	return "publisher"
}

func (pub *Publisher) Stop() error {
	cxt, cancel := context.WithTimeout(context.Background(), container.Mgr.Config.Server.ShutdownTimeout*time.Second)
	defer cancel()
	err := pub.server.Shutdown(cxt)
	if err != nil {
		container.Mgr.Logger.Printf("\"%s\" server stop error: %v\n", pub.GetName(), err)
	}
	return nil
}

func publisherHandler(c *gin.Context) {
	var params PublishMessage

	defer c.Request.Body.Close()

	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, errors.New("publisher/params_error", "invalid request payload"))
		return
	}

	container.Mgr.Dispatcher.Publish(params.Topics, params.Action, params.Body)
	container.Mgr.Dispatcher.NotifyNodeMessage(params.Topics, params.Action, params.Body)

	c.AbortWithStatus(http.StatusNoContent)
	return
}
