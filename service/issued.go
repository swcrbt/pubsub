package service

import (
	"go-issued-service/library/container"
	"go-issued-service/library/gracehttp"
	"go-issued-service/library/websocket"
	"net/http"
	"time"
)

type Issued struct {
	handler func(http.ResponseWriter, *http.Request)
}

func NewIssued() *Issued {
	return &Issued{
		handler: handler,
	}
}

func (iss *Issued) Run() error {
	http.HandleFunc("/issued", iss.handler)

	go func() {
		container.Mgr.Logger.Printf("Server Run At: \"%s\"\n", container.Mgr.Config.Server.Address)

		if err := gracehttp.ListenAndServe(container.Mgr.Config.Server.Address, nil); err != nil {
			container.Mgr.Logger.Println("Server error")
		}
	}()

	return nil
}

func (iss *Issued) GetName() string {
	return "issued"
}

func handler(w http.ResponseWriter, r *http.Request) {
	var (
		err  error
		conn *websocket.Connection
		data []byte
	)

	conn, err = websocket.New(w, r)
	if err != nil {
		container.Mgr.Logger.Println("websocket create failed: %v", err)
		return
	}

	go func() {
		for {
			if err = conn.WriteMessage([]byte("heartbeat")); err != nil {
				return
			}
			time.Sleep(1 * time.Second)
		}
	}()

	for {
		if data, err = conn.ReadMessage(); err != nil {
			goto ERR
		}
		if err = conn.WriteMessage(data); err != nil {
			goto ERR
		}
	}

ERR:
	conn.Close()
}
