package websocket

import (
	"errors"
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

type Connection struct {
	wsConnect *websocket.Conn
	inChan    chan []byte
	outChan   chan []byte
	closeChan chan byte

	delaySendPing chan byte
	mutex         sync.Mutex // 对closeChan关闭上锁
	isClosed      bool       // 防止closeChan被关闭多次
}

func New(w http.ResponseWriter, r *http.Request) (conn *Connection, err error) {
	var wsConn *websocket.Conn

	if wsConn, err = upgrader.Upgrade(w, r, nil); err != nil {
		return nil, err
	}

	conn = &Connection{
		wsConnect:     wsConn,
		inChan:        make(chan []byte, 1000),
		outChan:       make(chan []byte, 1000),
		delaySendPing: make(chan byte, 10),
		closeChan:     make(chan byte, 1),
	}

	// 启动读协程
	go conn.readLoop()
	// 启动写协程
	go conn.writeLoop()

	return conn, nil
}

func (conn *Connection) SetDeadline(readDeadline time.Duration, writeDeadline time.Duration) {
	_ = conn.wsConnect.SetReadDeadline(time.Now().Add(readDeadline))

	// 收到ping处理
	conn.wsConnect.SetPingHandler(func(message string) error {
		conn.delaySendPing <- 1
		err := conn.wsConnect.WriteControl(websocket.PongMessage, []byte(message), time.Now().Add(writeDeadline))
		if err == websocket.ErrCloseSent {
			return nil
		} else if e, ok := err.(net.Error); ok && e.Temporary() {
			return nil
		}
		return err
	})

	// 收到pong处理
	conn.wsConnect.SetPongHandler(func(message string) error {
		return conn.wsConnect.SetReadDeadline(time.Now().Add(readDeadline))
	})

	// 当没送到ping和text message时发送ping
	go func() {
		pingPeriod := (readDeadline * 9) / 10
		timer := time.NewTimer(pingPeriod)
		defer timer.Stop()

		for {
			select {
			case <-conn.delaySendPing:
				_ = conn.wsConnect.SetReadDeadline(time.Now().Add(readDeadline))
			case <-timer.C:
				if err := conn.wsConnect.WriteControl(websocket.PingMessage, nil, time.Now().Add(writeDeadline)); err != nil {
					return
				}
			}
			timer.Reset(pingPeriod)
		}
	}()
}

func (conn *Connection) ReadMessage() (data []byte, err error) {
	select {
	case data = <-conn.inChan:
		conn.delaySendPing <- 1
	case <-conn.closeChan:
		err = errors.New("connection is closeed")
	}
	return
}

func (conn *Connection) WriteMessage(data []byte) (err error) {
	select {
	case conn.outChan <- data:
	case <-conn.closeChan:
		err = errors.New("connection is closeed")
	}
	return
}

func (conn *Connection) Close() {
	// 线程安全，可多次调用
	conn.wsConnect.Close()
	// 利用标记，让closeChan只关闭一次
	conn.mutex.Lock()
	if !conn.isClosed {
		close(conn.closeChan)
		conn.isClosed = true
	}
	conn.mutex.Unlock()
}

// 内部实现
func (conn *Connection) readLoop() {
	var (
		data []byte
		err  error
	)

	for {
		if _, data, err = conn.wsConnect.ReadMessage(); err != nil {
			goto ERR
		}

		//阻塞在这里，等待inChan有空闲位置
		select {
		case conn.inChan <- data:
		case <-conn.closeChan: // closeChan 感知 conn断开
			goto ERR
		}
	}

ERR:
	conn.Close()
}

func (conn *Connection) writeLoop() {
	var (
		data []byte
		err  error
	)

	for {
		//阻塞在这里，等待outChan有数据
		select {
		case data = <-conn.outChan:
		case <-conn.closeChan:
			goto ERR
		}
		if err = conn.wsConnect.WriteMessage(websocket.TextMessage, data); err != nil {
			goto ERR
		}
	}

ERR:
	conn.Close()
}
