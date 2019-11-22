package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.orayer.com/golang/server"
	"gitlab.orayer.com/golang/issue/library/container"
	"gitlab.orayer.com/golang/issue/library/websocket"
	"gitlab.orayer.com/golang/issue/middleware"
	"net/http"
	"time"
)

type Issuer struct {
	handler func(c *gin.Context)
	server *server.HttpServer
}

func NewIssuer() *Issuer {
	return &Issuer{
		handler: issuerHandler,
	}
}

func (iss *Issuer) Run() error {

	//rev := gin.Default()
	//rev.GET("/subscribe", iss.handler)

	rev, err := server.Endless().NewHttpServer()
	if (err != nil) {
		return err
	}

	gin.SetMode(container.Mgr.Config.Server.Mode)
	rev.Router.Use(gin.Logger(), gin.Recovery(), middleware.IssueAuth())

	rev.Router.GET("/subscribe", iss.handler)
	rev.Port = container.Mgr.Config.Server.Issuer.Port

	go func() {
		container.Mgr.Logger.Printf("\"%s\" Server Run At: \"%d\"\n", iss.GetName(), container.Mgr.Config.Server.Issuer.Port)

		if err := rev.Start(); err != nil {
			container.Mgr.Logger.Printf("\"%s\" Server error: %v\n", iss.GetName(), err)
		}

		/*if err := rev.Run(":" + strconv.Itoa(container.Mgr.Config.Server.Issuer.Port)); err != nil {
			container.Mgr.Logger.Printf("\"%s\" Server error: %v\n", iss.GetName(), err)
		}*/
	}()

	iss.server = rev

	return nil
}

func (iss *Issuer) GetName() string {
	return "issuer"
}

func (rec *Issuer) Stop() error  {
	if rec.server != nil {
		return rec.server.Shutdown()
	}
	return nil
}

func issuerHandler(c *gin.Context) {
	var (
		dataChan     <-chan []byte
		err          error
		conn         *websocket.Connection
		timeoutCount int = 0
		index        int = 0
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

	dataChan, index = container.Mgr.Dispatcher.Subscribe(auth.(middleware.StorageRecord).Action, auth.(middleware.StorageRecord).UniqId)

	// 发送心跳包
	go func() {
		for {
			if timeoutCount >= container.Mgr.Config.Server.Issuer.HeartbeatTimeout {
				conn.Close()
				container.Mgr.Dispatcher.UnSubscribe(auth.(middleware.StorageRecord).Action, auth.(middleware.StorageRecord).UniqId, index)
				return
			}

			if err = conn.WriteMessage([]byte("_heartbeat_")); err != nil {
				fmt.Println(err)
				container.Mgr.Dispatcher.UnSubscribe(auth.(middleware.StorageRecord).Action, auth.(middleware.StorageRecord).UniqId, index)
				return
			}

			timeoutCount++
			time.Sleep(container.Mgr.Config.Server.Issuer.HeartbeatInterval * time.Second)
		}
	}()

	// 检测心跳包回应
	go func() {
		for {
			data, err := conn.ReadMessage()
			if err != nil {
				return
			}

			if string(data) == "NOP" {
				timeoutCount = 0;
			}
		}
	}()

	// 下发由调度器来的数据
	for {
		if err = conn.WriteMessage(<-dataChan); err != nil {
			break
		}
	}

	conn.Close()
}
