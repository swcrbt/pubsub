package dispatcher

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

type PublishRecord struct {
	ID     uint        `json:"id"`
	Action string      `json:"action" binding:"required"`
	Body   interface{} `json:"body"`
}

type Dispatcher struct {
	pool map[string]map[uint]chan []byte
	back map[uint]map[string]map[uint]chan *PublishRecord

	mid    uint
	cid    uint
	mMutex sync.Mutex
	cMutex sync.Mutex
}

func New() *Dispatcher {
	return &Dispatcher{
		pool: make(map[string]map[uint]chan []byte, 0xffff),
		back: make(map[uint]map[string]map[uint]chan *PublishRecord, 0xffff),
		mid:  0,
		cid:  0,
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

	mid := dis.GetMId()
	resp := PublishRecord{ID: mid, Action: action, Body: body}
	respData, err := json.Marshal(resp)
	if err != nil {
		return result
	}

	// 初始化回调池
	if _, ok := dis.back[mid]; !ok {
		dis.back[mid] = map[string]map[uint]chan *PublishRecord{}
	}

	// 下发数据
	for _, topic := range topics {
		if listeners, ok := dis.pool[topic]; ok {
			for cid, listener := range listeners {
				if listener != nil {
					if _, ok := dis.back[mid][topic]; !ok {
						dis.back[mid][topic] = map[uint]chan *PublishRecord{}
					}

					callbacker := make(chan *PublishRecord, 1)
					dis.back[mid][topic][cid] = callbacker

					listener <- respData
				}
			}
		}
	}

	// 等待数据回调
	for topic, callbackers := range dis.back[mid] {
		for cid, callbacker := range callbackers {
			timer := time.NewTimer(time.Second * 10)
			select {
			case data := <-callbacker:
				fmt.Printf("%s: %v -- %v\n", topic, cid, data)
			case <-timer.C:
				fmt.Printf("timeout C %s: %v\n", topic, cid)
				break
			}
		}
	}

	return result
}

func (dis *Dispatcher) Feedback(topic string, cid uint, data *PublishRecord) {
	if _, ok := dis.back[data.ID]; !ok {
		return
	}

	if _, ok := dis.back[data.ID][topic]; !ok {
		return
	}

	if callbacker, ok := dis.back[data.ID][topic][cid]; ok {
		callbacker <- data
		close(callbacker)
	}
}

func (dis *Dispatcher) Subscribe(topic string) (<-chan []byte, uint) {
	listener := make(chan []byte, 100)
	cid := dis.GetCId()

	if _, ok := dis.pool[topic]; !ok {
		dis.pool[topic] = map[uint]chan []byte{}
	}

	dis.pool[topic][cid] = listener

	fmt.Printf("pool: %v\n", dis.pool)

	return listener, cid
}

func (dis *Dispatcher) UnSubscribe(topic string, cid uint) {
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

func (dis *Dispatcher) GetMId() uint {
	dis.mMutex.Lock()

	id := dis.mid + 1
	dis.mid = id

	dis.mMutex.Unlock()

	return id
}

func (dis *Dispatcher) GetCId() uint {
	dis.cMutex.Lock()

	id := dis.cid + 1
	dis.cid = id

	dis.cMutex.Unlock()

	return id
}
