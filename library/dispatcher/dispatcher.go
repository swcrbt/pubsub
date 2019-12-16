package dispatcher

import (
	"encoding/json"
	"github.com/google/uuid"
	"sync"
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
	subers sync.Map // map[string]map[string]chan []byte
	// msgbacks sync.Map // map[string]*MsgBackRecord
}

func New() *Dispatcher {
	return &Dispatcher{}
}

func (dis *Dispatcher) Publish(topics []string, action string, body interface{}) {
	// result := &sync.Map{} // map[string]map[string]bool{}

	// if len(topics) == 0 {
	// 	return result
	// }

	// wg := &sync.WaitGroup{}

	// 下发数据
	for _, topic := range topics {

		if listeners, ok := dis.subers.Load(topic); ok {
			listeners.(*sync.Map).Range(func(cid, listener interface{}) bool {
				if listener == nil {
					return true
				}

				mid := uuid.New().String()
				resp := PublishRecord{ID: mid, Action: action, Body: body}
				respData, err := json.Marshal(resp)
				if err != nil {
					return true
				}

				/*callbacker := make(chan *PublishRecord, 1)
				dis.msgbacks.Store(mid, &MsgBackRecord{
					Topic: topic,
					Chan:  callbacker,
				})*/

				listener.(chan []byte) <- respData

				/*wg.Add(1)
				go func(topic string, cid string, callbacker chan *PublishRecord) {
					timer := time.NewTimer(time.Second * 5)
					defer func() {
						wg.Done()
						timer.Stop()
					}()
					res, _ := result.LoadOrStore(topic, &sync.Map{})
					select {
					case <-callbacker:
						res.(*sync.Map).Store(cid, true)
					case <-timer.C:
						res.(*sync.Map).Store(cid, false)
					}
				}(topic, cid.(string), callbacker)*/

				return true
			})
		}
	}

	// wg.Wait()

	// return result
}

func (dis *Dispatcher) Subscribe(topic string, cid string, listener chan []byte) {
	if topic == "" || cid == "" {
		return
	}

	suber, _ := dis.subers.LoadOrStore(topic, &sync.Map{})
	suber.(*sync.Map).Store(cid, listener)
}

func (dis *Dispatcher) UnSubscribe(topic string, cid string) {
	if topic == "" || cid == "" {
		return
	}

	if subers, ok := dis.subers.Load(topic); ok {
		subers.(*sync.Map).Delete(cid)
	}
}

/*func (dis *Dispatcher) Feedback(cid string, data *PublishRecord) {
	if callbacker, ok := dis.msgbacks.Load(data.ID); ok {
		callbacker.(*MsgBackRecord).Chan <- data
		dis.msgbacks.Delete(data.ID)
	}
}*/

func (dis *Dispatcher) Destroy(cid string, listener chan []byte) {
	dis.subers.Range(func(topic, callbackers interface{}) bool {

		callbackers.(*sync.Map).Range(func(id, _ interface{}) bool {
			if id.(string) == cid {
				callbackers.(*sync.Map).Store(cid, nil)
				callbackers.(*sync.Map).Delete(cid)
			}
			return true
		})

		return true
	})
}
