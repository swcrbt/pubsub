package service

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.orayer.com/golang/issue/library/container"
	"gitlab.orayer.com/golang/issue/library/dispatcher"
	"gitlab.orayer.com/golang/issue/library/websocket"
	"gitlab.orayer.com/golang/issue/middleware"
	"gitlab.orayer.com/golang/server"
	"net/http"
	"time"
)

type Subscriber struct {
	handler func(c *gin.Context)
	server  *server.HttpServer
}

func NewSubscriber() *Subscriber {
	return &Subscriber{
		handler: subscriberHandler,
	}
}

func (iss *Subscriber) Run() error {
	rev := server.NewHttpServer()

	rev.Router.Use(gin.Logger(), gin.Recovery(), middleware.SubAuth())

	rev.Router.GET("/subscribe", iss.handler)
	rev.Port = container.Mgr.Config.Server.Subscriber.Port

	go func() {
		container.Mgr.Logger.Printf("\"%s\" Server Run At: \"%d\"\n", iss.GetName(), container.Mgr.Config.Server.Subscriber.Port)

		if err := rev.Start(); err != nil {
			container.Mgr.Logger.Printf("\"%s\" Server error: %v\n", iss.GetName(), err)
		}
	}()

	iss.server = rev

	return nil
}

func (iss *Subscriber) GetName() string {
	return "subscriber"
}

func (rec *Subscriber) Stop() error {
	if rec.server != nil {
		return rec.server.Shutdown()
	}
	return nil
}

func subscriberHandler(c *gin.Context) {
	var (
		err  error
		conn *websocket.Connection
	)

	auth, ok := c.Get(middleware.AuthInfoKey)
	if !ok {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	conn, err = websocket.New(c.Writer, c.Request)
	if err != nil {
		container.Mgr.Logger.Printf("websocket create failed: %v\n", err)
		return
	}

	conn.SetDeadline(
		time.Second*container.Mgr.Config.Server.Subscriber.ReadDeadline,
		time.Second*container.Mgr.Config.Server.Subscriber.WriteDeadline,
	)

	cid := auth.(middleware.StorageRecord).ChannelID
	topics := auth.(middleware.StorageRecord).Topics
	channel := make(chan []byte, 10)

	fmt.Printf("cid:\"%v\" channel: %v\n", cid, channel)

	for _, topic := range topics {
		container.Mgr.Dispatcher.Subscribe(topic, cid, channel)
	}

	// 回应
	go func() {
		for {
			data, err := conn.ReadMessage()
			if err != nil {
				container.Mgr.Logger.Printf("cid:\"%v\" read message failed: %v\n", cid, err)
				return
			}

			fmt.Printf("cid:\"%v\" read message: %s\n", cid, string(data))

			var resp *dispatcher.PublishRecord
			err = json.Unmarshal(data, &resp)
			if err != nil {
				continue
			}

			if resp.Action == "subscribe" {
				if body, ok := resp.Body.(map[string]interface{}); ok {
					container.Mgr.Dispatcher.Subscribe(body["topic"].(string), cid, channel)
				}
				continue
			}

			if resp.Action == "unsubscribe" {
				if body, ok := resp.Body.(map[string]interface{}); ok {
					container.Mgr.Dispatcher.UnSubscribe(body["topic"].(string), cid)
				}
				continue
			}

			container.Mgr.Dispatcher.Feedback(cid, resp)
		}
	}()

	// 下发由调度器来的数据
	for {
		data, ok := <-channel
		if (!ok) {
			break
		}

		if err = conn.WriteMessage(data); err != nil {
			container.Mgr.Logger.Printf("cid:\"%v\" write data failed: %v\n", cid, err)
			break
		}
	}

	defer func() {
		container.Mgr.Dispatcher.Destroy(cid, channel)
		conn.Close()
	}()
}
