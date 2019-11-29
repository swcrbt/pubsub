package service

import (
	"encoding/json"
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

	rev.Router.Use(gin.Logger(), gin.Recovery(), middleware.IssueAuth())

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
		err          error
		conn         *websocket.Connection
		timeoutCount int = 0
	)

	auth, ok := c.Get(middleware.AuthInfoKey)
	if !ok {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	conn, err = websocket.New(c.Writer, c.Request)
	if err != nil {
		container.Mgr.Logger.Printf("websocket create failed: %v\n", err)
		return
	}

	topic := auth.(middleware.StorageRecord).Topic
	channel, cid:= container.Mgr.Dispatcher.Subscribe(topic)

	// 发送心跳包
	go func() {
		for {
			container.Mgr.Logger.Printf("topic:\"%s\" cid:\"%v\" count:%d\n", topic, cid, timeoutCount)

			if timeoutCount >= container.Mgr.Config.Server.Subscriber.HeartbeatTimeout {
				container.Mgr.Logger.Printf("topic:\"%s\" cid:\"%v\" timeout\n", topic, cid)
				container.Mgr.Dispatcher.UnSubscribe(topic, cid)
				conn.Close()
				return
			}

			hb, _ := json.Marshal(dispatcher.PublishRecord{ID: "0", Action: "heartbeat"})
			if err = conn.WriteMessage(hb); err != nil {
				container.Mgr.Logger.Printf("topic:\"%s\" cid:\"%v\" write heartbeat failed: %v\n", topic, cid, err)
				container.Mgr.Dispatcher.UnSubscribe(topic, cid)
				return
			}

			timeoutCount++
			time.Sleep(container.Mgr.Config.Server.Subscriber.HeartbeatInterval * time.Second)
		}
	}()

	// 回应
	go func() {
		for {
			data, err := conn.ReadMessage()
			if err != nil {
				container.Mgr.Logger.Printf("topic:\"%s\" cid:\"%v\" read message failed: %v\n", topic, cid, err)
				return
			}

			container.Mgr.Logger.Printf("topic:\"%s\" cid:\"%v\" read message: %s\n", topic, cid, string(data))

			var resp *dispatcher.PublishRecord
			err = json.Unmarshal(data, &resp)
			if err != nil {
				continue
			}

			if resp.ID == "0" || resp.ID == "" {
				timeoutCount = 0;
				continue
			}

			container.Mgr.Dispatcher.Feedback(topic, cid, resp)
		}
	}()

	// 下发由调度器来的数据
	for {
		data, ok := <-channel
		if (!ok) {
			break
		}

		if err = conn.WriteMessage(data); err != nil {
			container.Mgr.Logger.Printf("topic:\"%s\" cid:\"%v\" release data failed: %v\n", topic, cid, err)
			break
		}
	}

	conn.Close()
}
