package dispatcher

import (
	"github.com/google/uuid"
	"time"
)

type Channel struct {
	id         string
	topicName  string
	clientID   string
	createTime time.Time
}

func NewChannel(topicName string, clientID string) *Channel {
	return &Channel{
		id:         uuid.New().String(),
		topicName:  topicName,
		clientID:   clientID,
		createTime: time.Now(),
	}
}
