package dispatcher

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
	"github.com/google/uuid"
)

type PublishRecord struct {
	ID     string      `json:"id"`
	Action string      `json:"action" binding:"required"`
	Body   interface{} `json:"body"`
}

type Dispatcher struct {
	pool map[string]map[string]chan []byte
	back map[string]map[string]map[string]chan *PublishRecord
}

func New() *Dispatcher {
	return &Dispatcher{
		pool: make(map[string]map[string]chan []byte, 0xffff),
		back: make(map[string]map[string]map[string]chan *PublishRecord, 0xffff),
	}
}

func (dis *Dispatcher) Publish(topics []string, action string, body interface{}) map[string]uint32 {
	result := map[string]uint32{}

	if len(topics) == 0 {
		return result
	}

	for _, topic := range topics {
		result[topic] = 0
	}

	mid := uuid.New().String()
	resp := PublishRecord{ID: mid, Action: action, Body: body}
	respData, err := json.Marshal(resp)
	if err != nil {
		return result
	}

	// 初始化回调池
	if _, ok := dis.back[mid]; !ok {
		dis.back[mid] = map[string]map[string]chan *PublishRecord{}
	}

	// 下发数据
	for _, topic := range topics {
		if listeners, ok := dis.pool[topic]; ok {
			for cid, listener := range listeners {
				if listener != nil {
					if _, ok := dis.back[mid][topic]; !ok {
						dis.back[mid][topic] = map[string]chan *PublishRecord{}
					}

					callbacker := make(chan *PublishRecord, 1)
					dis.back[mid][topic][cid] = callbacker

					listener <- respData
				}
			}
		}
	}

	wg := &sync.WaitGroup{}

	// 等待数据回调
	for topic, callbackers := range dis.back[mid] {
		for cid, callbacker := range callbackers {
			wg.Add(1)
			go func(topic string, cid string, callbacker chan *PublishRecord) {
				timer := time.NewTimer(time.Second * 5)
				select {
				case data := <-callbacker:
					result[topic] ++
					fmt.Printf("back %s[%v] -- %v\n", topic, cid, data)
				case <-timer.C:
					fmt.Printf("back timeout %s[%v] \n", topic, cid)
				}
				wg.Done()
			}(topic, cid, callbacker)
		}
	}

	wg.Wait()

	return result
}

func (dis *Dispatcher) Feedback(topic string, cid string, data *PublishRecord) {
	if _, ok := dis.back[data.ID]; !ok {
		return
	}

	if _, ok := dis.back[data.ID][topic]; !ok {
		return
	}

	fmt.Printf("feedback [%v] %s[%v] -- %v\n", data.ID, topic, cid, data)

	if callbacker, ok := dis.back[data.ID][topic][cid]; ok {
		callbacker <- data
		delete(dis.back[data.ID][topic], cid)
	}
}

func (dis *Dispatcher) Subscribe(topic string) (<-chan []byte, string) {
	listener := make(chan []byte, 100)
	cid := uuid.New().String()

	if _, ok := dis.pool[topic]; !ok {
		dis.pool[topic] = map[string]chan []byte{}
	}

	dis.pool[topic][cid] = listener

	fmt.Printf("pool: %v\n", dis.pool)

	return listener, cid
}

func (dis *Dispatcher) UnSubscribe(topic string, cid string) {
	if _, ok := dis.pool[topic]; !ok {
		return
	}

	if listener, ok := dis.pool[topic][cid]; !ok || listener == nil {
		return
	}

	//close(dis.pool[topic][uid])
	//dis.pool[topic][uid] = nil
	delete(dis.pool[topic], cid)

	fmt.Printf("pool: %v\n", dis.pool)
}