package dispatcher

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"net"
	"net/http"
	"sync"
	"time"
)

var upgrader = websocket.Upgrader{
	// 允许跨域
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	sync.RWMutex

	ID        string
	wsConnect *websocket.Conn

	ConnectTime time.Time

	InChan  chan []byte
	OutChan chan []byte

	isExit  bool
	isDelaySendPing bool
	exitChan        chan byte
}

func NewClient(ctx *gin.Context) (*Client, error) {
	wsConn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil);
	if err != nil {
		return nil, err
	}

	client := &Client{
		ID:              uuid.New().String(),
		wsConnect:       wsConn,
		ConnectTime:	 time.Now(),
		InChan:          make(chan []byte, 1000),
		OutChan:         make(chan []byte, 1000),
		isDelaySendPing: false,
		exitChan:        make(chan byte, 1),
	}

	return client, nil
}

func (c *Client) SetDeadline(readDeadline time.Duration, writeDeadline time.Duration) {
	_ = c.wsConnect.SetReadDeadline(time.Now().Add(readDeadline))

	// 收到ping处理
	c.wsConnect.SetPingHandler(func(message string) error {
		c.checkOrSetDelayPing(true)
		err := c.wsConnect.WriteControl(websocket.PongMessage, []byte(message), time.Now().Add(writeDeadline))
		if err == websocket.ErrCloseSent {
			return nil
		} else if e, ok := err.(net.Error); ok && e.Temporary() {
			return nil
		}
		return err
	})

	// 收到pong处理
	c.wsConnect.SetPongHandler(func(message string) error {
		return c.wsConnect.SetReadDeadline(time.Now().Add(readDeadline))
	})

	// 当没送到ping和text message时发送ping
	go func() {
		pingPeriod := (readDeadline * 9) / 10
		ticker := time.NewTicker(pingPeriod)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if c.checkOrSetDelayPing(false) {
					_ = c.wsConnect.SetReadDeadline(time.Now().Add(readDeadline))
					continue
				}
				if err := c.wsConnect.WriteControl(websocket.PingMessage, nil, time.Now().Add(writeDeadline)); err != nil {
					return
				}
			case <-c.exitChan:
				return
			}
		}
	}()
}

func (c *Client) ReadMessage() (data []byte, err error) {
	select {
	case data = <-c.InChan:
		c.checkOrSetDelayPing(true)
	case <-c.exitChan:
		err = errors.New("connection is closeed")
	}
	return
}

func (c *Client) WriteMessage(data []byte) (err error) {
	select {
	case c.OutChan <- data:
	case <-c.exitChan:
		err = errors.New("connection is closeed")
	}
	return
}

func (c *Client) Start() {
	// 启动读协程
	go c.readLoop()
	// 启动写协程
	go c.writeLoop()
}

func (c *Client) GetIsExit() bool {
	c.RLock()
	defer c.RUnlock()
	return c.isExit
}

func (c *Client) Exit() {
	// 线程安全，可多次调用
	c.wsConnect.Close()
	// 利用标记，让closeChan只关闭一次
	c.Lock()
	if !c.isExit {
		close(c.exitChan)
		c.isExit = true
	}
	c.Unlock()
}

// 内部实现
func (c *Client) readLoop() {
	var (
		data []byte
		err  error
	)

	for {
		if _, data, err = c.wsConnect.ReadMessage(); err != nil {
			goto ERR
		}

		//阻塞在这里，等待inChan有空闲位置
		select {
		case c.InChan <- data:
		case <-c.exitChan: // closeChan 感知 conn断开
			goto ERR
		}
	}

ERR:
	c.Exit()
}

func (c *Client) writeLoop() {
	var (
		data []byte
		err  error
	)

	for {
		//阻塞在这里，等待outChan有数据
		select {
		case data = <-c.OutChan:
		case <-c.exitChan:
			goto ERR
		}
		if err = c.wsConnect.WriteMessage(websocket.TextMessage, data); err != nil {
			goto ERR
		}
	}

ERR:
	c.Exit()
}

func (c *Client) checkOrSetDelayPing(isDelay bool) bool {
	c.Lock()
	defer c.Unlock()

	result := c.isDelaySendPing
	c.isDelaySendPing = isDelay

	return result
}
