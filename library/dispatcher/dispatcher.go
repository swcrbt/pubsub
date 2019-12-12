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
	subers   map[string]map[string]chan []byte
	msgbacks map[string]*MsgBackRecord

	sublock sync.Mutex
	msglock sync.Mutex
}

func New() *Dispatcher {
	return &Dispatcher{
		subers:   make(map[string]map[string]chan []byte, 0xffff),
		msgbacks: map[string]*MsgBackRecord{},
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

		if listeners, ok := dis.subers[topic]; ok {
			dis.msglock.Lock()
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
				dis.msgbacks[mid] = &MsgBackRecord{
					Topic: topic,
					Chan:  callbacker,
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
			dis.msglock.Unlock()
		}
	}

	fmt.Printf("callbackers back:%v \n", dis.msgbacks)

	wg.Wait()

	return result
}

func (dis *Dispatcher) Subscribe(topic string, cid string, listener chan []byte) {
	if topic == "" || cid == "" {
		return
	}

	dis.sublock.Lock()
	defer dis.sublock.Unlock()

	if _, ok := dis.subers[topic]; !ok {
		dis.subers[topic] = map[string]chan []byte{}
	}

	/*if _, ok := dis.subers[topic][cid]; ok {
		dis.UnSubscribe(topic, cid)
	}*/

	dis.subers[topic][cid] = listener

	fmt.Printf("sub subers: %v\n", dis.subers)
}

func (dis *Dispatcher) UnSubscribe(topic string, cid string) {
	if topic == "" || cid == "" {
		return
	}

	if _, ok := dis.subers[topic]; !ok {
		return
	}

	if listener, ok := dis.subers[topic][cid]; !ok || listener == nil {
		return
	}

	dis.sublock.Lock()
	defer dis.sublock.Unlock()

	delete(dis.subers[topic], cid)

	fmt.Printf("unsub subers: %v\n", dis.subers)
}

func (dis *Dispatcher) Feedback(cid string, data *PublishRecord) {
	if _, ok := dis.msgbacks[data.ID]; !ok {
		return
	}

	dis.msglock.Lock()
	defer dis.msglock.Unlock()

	fmt.Printf("feedback %s[%v] -- %v\n", data.ID, cid, data)

	if callbacker, ok := dis.msgbacks[data.ID]; ok {
		callbacker.Chan <- data
		delete(dis.msgbacks, data.ID)
	}
}

func (dis *Dispatcher) Destroy(cid string, listener chan []byte) {
	dis.sublock.Lock()
	defer dis.sublock.Unlock()

	for topic, callbackers := range dis.subers {
		for id := range callbackers {
			if id == cid {
				delete(dis.subers[topic], id)
			}
		}
	}

	if listener != nil {
		close(listener)
	}
}
