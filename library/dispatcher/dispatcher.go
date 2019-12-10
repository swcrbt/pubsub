package dispatcher

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"sync"
	"time"
)

type PublishRecord struct {
	ID     string      `json:"id"`
	Action string      `json:"action" binding:"required"`
	Body   interface{} `json:"body"`
}

type MsgBackRecord struct {
	Topic string
	Chan  chan *PublishRecord
}

type Dispatcher struct {
	pool    map[string]map[string]chan []byte
	msgback map[string]*MsgBackRecord
}

func New() *Dispatcher {
	return &Dispatcher{
		pool:    make(map[string]map[string]chan []byte, 0xffff),
		msgback: map[string]*MsgBackRecord{},
	}
}

func (dis *Dispatcher) Publish(topics []string, action string, body interface{}) map[string]map[string]bool {
	result := map[string]map[string]bool{}

	if len(topics) == 0 {
		return result
	}

	wg := &sync.WaitGroup{}

	// 下发数据
	for _, topic := range topics {

		if _, ok := result[topic]; !ok {
			result[topic] = map[string]bool{}
		}

		if listeners, ok := dis.pool[topic]; ok {
			for cid, listener := range listeners {
				if listener == nil {
					continue
				}

				mid := uuid.New().String()
				resp := PublishRecord{ID: mid, Action: action, Body: body}
				respData, err := json.Marshal(resp)
				if err != nil {
					continue
				}

				callbacker := make(chan *PublishRecord, 1)
				dis.msgback[mid] = &MsgBackRecord{
					Topic: topic,
					Chan: callbacker,
				}

				listener <- respData

				wg.Add(1)
				go func(topic string, cid string, callbacker chan *PublishRecord) {
					timer := time.NewTimer(time.Second * 5)
					select {
					case data := <-callbacker:
						result[topic][cid] = true
						fmt.Printf("back %s[%v] -- %v\n", topic, cid, data)
					case <-timer.C:
						result[topic][cid] = false
						fmt.Printf("back timeout %s[%v] \n", topic, cid)
					}
					wg.Done()
				}(topic, cid, callbacker)
			}
		}
	}

	fmt.Printf("callbackers back:%v \n", dis.msgback)

	wg.Wait()

	return result
}

func (dis *Dispatcher) Subscribe(topic string, cid string, listener chan []byte) {
	if topic == "" || cid == "" {
		return
	}

	if _, ok := dis.pool[topic]; !ok {
		dis.pool[topic] = map[string]chan []byte{}
	}

	/*if _, ok := dis.pool[topic][cid]; ok {
		dis.UnSubscribe(topic, cid)
	}*/

	dis.pool[topic][cid] = listener

	fmt.Printf("sub pool: %v\n", dis.pool)
}

func (dis *Dispatcher) UnSubscribe(topic string, cid string) {
	if topic == "" || cid == "" {
		return
	}

	if _, ok := dis.pool[topic]; !ok {
		return
	}

	if listener, ok := dis.pool[topic][cid]; !ok || listener == nil {
		return
	}

	delete(dis.pool[topic], cid)

	fmt.Printf("unsub pool: %v\n", dis.pool)
}

func (dis *Dispatcher) Feedback(cid string, data *PublishRecord) {
	if _, ok := dis.msgback[data.ID]; !ok {
		return
	}

	fmt.Printf("feedback %s[%v] -- %v\n", data.ID, cid, data)

	if callbacker, ok := dis.msgback[data.ID]; ok {
		callbacker.Chan <- data
		delete(dis.msgback, data.ID)
	}
}

func (dis *Dispatcher) Destroy(cid string, listener chan []byte) {
	for topic, callbackers := range dis.pool {
		for id := range callbackers {
			if id == cid {
				delete(dis.pool[topic], id)
			}
		}
	}

	if listener != nil {
		close(listener)
		fmt.Println(listener)
		//listener = nil
	}
}
