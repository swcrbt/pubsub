package dispatcher

import (
	"sync"
)

type Topic struct {
	sync.RWMutex

	name             string
	ClientChannelMap map[string]*Channel
}

func NewTopic(topicName string) *Topic {
	return &Topic{
		name:             topicName,
		ClientChannelMap: map[string]*Channel{},
	}
}

func (t *Topic) AddClientChannel(clientID string) {
	t.Lock()
	t.ClientChannelMap[clientID] = NewChannel(t.name, clientID)
	t.Unlock()
}

func (t *Topic) RemoveClientChannel(clientID string) {
	t.Lock()
	delete(t.ClientChannelMap, clientID)
	t.Unlock()
}