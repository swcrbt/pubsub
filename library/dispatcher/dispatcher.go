package dispatcher

import (
	"encoding/json"
	"fmt"
)

type Dispatcher struct {
	pool map[string][]chan []byte
}

func New() *Dispatcher {
	return &Dispatcher{
		pool: map[string][]chan []byte{},
	}
}

func (dis *Dispatcher) Push(action string, uids []string, data interface{}) {
	d, err := json.Marshal(data)
	if err == nil {
		for _, uid := range uids {
			key := action + "_" + uid
			fmt.Println(key, dis.pool)
			if listeners, ok := dis.pool[key]; ok {
				for _, listener := range listeners {
					if listener != nil {
						listener <- d
					}
				}
			}
		}
	}
}

func (dis *Dispatcher) Subscribe(action string, uid string) (<-chan []byte, int) {
	listener := make(chan []byte, 100)
	key := action + "_" + uid
	index := 0

	if _, ok := dis.pool[key]; ok {
		index = len(dis.pool[key])
		dis.pool[key] = append(dis.pool[key], listener)
	} else {
		dis.pool[key] = []chan []byte{listener}
	}

	x := dis.pool[key]

	fmt.Printf("key=%s len=%d cap=%d slice=%v\n", key, len(x), cap(x), x)

	return listener, index
}

func (dis *Dispatcher) UnSubscribe(action string, uid string, index int) {
	key := action + "_" + uid

	if _, ok := dis.pool[key]; ok {
		count := len(dis.pool[key])
		noNilCount := 0

		dis.pool[key][index] = nil
		for k, v := range dis.pool[key] {
			if v != nil {
				noNilCount++
			}

			if noNilCount > 0 {
				break
			}

			if (k == count-1) {
				delete(dis.pool, key)
			}
		}
	}
}
