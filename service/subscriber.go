package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.orayer.com/golang/pubsub/dispatcher"
	"gitlab.orayer.com/golang/pubsub/library/container"
	"gitlab.orayer.com/golang/pubsub/middleware"
	"net/http"
	"time"
)

type Subscriber struct {
	handler  func(c *gin.Context)
	server   *http.Server
}

func NewSubscriber() *Subscriber {
	return &Subscriber{
		handler:  subscriberHandler,
	}
}

func (sub *Subscriber) Run() error {
	router := gin.Default()
	router.Use(middleware.SubAuth())
	router.GET("/subscribe", sub.handler)

	rev := &http.Server{
		Addr:    container.Mgr.Config.Server.Subscriber.Address,
		Handler: router,
	}

	go func() {
		container.Mgr.Logger.Printf("\"%s\" server run at: \"%s\"\n", sub.GetName(), rev.Addr)

		if err := rev.ListenAndServe(); err != nil {
			container.Mgr.Logger.Printf("\"%s\" server stop at \"%s\": %v\n", sub.GetName(), rev.Addr, err)
		}
	}()

	sub.server = rev

	return nil
}

func (sub *Subscriber) GetName() string {
	return "subscriber"
}

func (sub *Subscriber) Stop() error {
	cxt, cancel := context.WithTimeout(context.Background(), container.Mgr.Config.Server.ShutdownTimeout*time.Second)
	defer cancel()
	err := sub.server.Shutdown(cxt)
	if err != nil {
		container.Mgr.Logger.Printf("\"%s\" Server stop error: %v\n", sub.GetName(), err)
	}
	return nil
}

func subscriberHandler(c *gin.Context) {
	auth, ok := c.Get(middleware.AuthInfoKey)
	if !ok {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	id := auth.(middleware.StorageRecord).ID
	topics := auth.(middleware.StorageRecord).Topics

	client, err := dispatcher.NewClient(c)
	if err != nil {
		container.Mgr.Logger.Printf("websocket create failed: %v\n", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	client.SetDeadline(
		time.Second*container.Mgr.Config.Server.Subscriber.ReadDeadline,
		time.Second*container.Mgr.Config.Server.Subscriber.WriteDeadline,
	)

	client.Start()

	container.Mgr.Dispatcher.AddClient(id, client)
	container.Mgr.Logger.Printf("cid: \"%v\" connecttime: %v \n", client.ID, client.ConnectTime);

	defer func() {
		container.Mgr.Dispatcher.RemoveClient(client.ID)
		client.Exit()
	}()

	for _, topic := range topics {
		container.Mgr.Dispatcher.Subscribe(topic, client.ID)
	}

	for {
		data, err := client.ReadMessage()
		if err != nil {
			container.Mgr.Logger.Printf("cid:\"%v\" read message failed: %v\n", client.ID, err)
			break
		}

		fmt.Printf("cid:\"%v\" read message: %s\n", client.ID, string(data))

		var resp *dispatcher.Message
		err = json.Unmarshal(data, &resp)
		if err != nil {
			continue
		}

		switch resp.Action {
		// 订阅
		case "sub":
			if body, ok := resp.Body.(map[string]interface{}); ok {
				container.Mgr.Dispatcher.Subscribe(body["topic"].(string), client.ID)
			}
		// 取订
		case "unsub":
			if body, ok := resp.Body.(map[string]interface{}); ok {
				container.Mgr.Dispatcher.UnSubscribe(body["topic"].(string), client.ID)
			}
		// 成功收到消息，后续可做重试
		case "fin":
		}
	}
}
