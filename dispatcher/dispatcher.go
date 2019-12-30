package dispatcher

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"gitlab.orayer.com/golang/pubsub/protos"
	"google.golang.org/grpc"
	"sync"
	"time"
)

type Message struct {
	ID     string      `json:"id"`
	Action string      `json:"action" binding:"required"`
	Body   interface{} `json:"body"`
}

type Dispatcher struct {
	sync.RWMutex

	createTime time.Time

	topicMap      map[string]*Topic
	clientMap     map[string]*Client
	nodeServerMap map[string]*Node
	subClientMap  map[string]string
}

func New() *Dispatcher {
	return &Dispatcher{
		createTime:    time.Now(),
		topicMap:      make(map[string]*Topic),
		clientMap:     make(map[string]*Client),
		nodeServerMap: make(map[string]*Node),
		subClientMap:  map[string]string{},
	}
}

func (dis *Dispatcher) SetNodeServer(nodeServer map[string]*Node) {
	dis.Lock()
	dis.nodeServerMap = nodeServer
	dis.Unlock()
}

func (dis *Dispatcher) AddClient(subID string, client *Client) {
	dis.Lock()
	if subID != "" {
		if clientID, ok := dis.subClientMap[subID]; ok {
			if cli, ok := dis.clientMap[clientID]; ok {
				cli.Exit()
			}
		}
		dis.subClientMap[subID] = client.ID
	}
	dis.clientMap[client.ID] = client
	dis.Unlock()
}

func (dis *Dispatcher) RemoveClient(clientID string) {
	dis.Lock()
	defer dis.Unlock()
	if _, ok := dis.clientMap[clientID]; !ok {
		return
	}
	delete(dis.clientMap, clientID)
}

func (dis *Dispatcher) Publish(topicNames []string, action string, body interface{}) {
	dis.RLock()
	defer dis.RUnlock()

	for _, topicName := range topicNames {
		topic, ok := dis.topicMap[topicName]
		if !ok {
			continue
		}

		topic.RLock()
		for clientID := range topic.ClientChannelMap {
			client, ok := dis.clientMap[clientID]
			if !ok || client.GetIsExit() {
				topic.RUnlock()
				topic.RemoveClientChannel(clientID)
				topic.RLock()
				continue
			}

			msg := Message{ID: uuid.New().String(), Action: action, Body: body}
			data, err := json.Marshal(msg)
			if err != nil {
				continue
			}

			_ = client.WriteMessage(data)
		}
		topic.RUnlock()
	}
}

func (dis *Dispatcher) Subscribe(topicName string, clientID string) {
	if topicName == "" || clientID == "" {
		return
	}

	dis.RLock()
	topic, ok := dis.topicMap[topicName]
	dis.RUnlock()

	if !ok {
		topic = NewTopic(topicName)
		dis.Lock()
		dis.topicMap[topicName] = topic
		dis.Unlock()
	}

	topic.AddClientChannel(clientID)
}

func (dis *Dispatcher) UnSubscribe(topicName string, clientID string) {
	if topicName == "" || clientID == "" {
		return
	}

	dis.RLock()
	if topic, ok := dis.topicMap[topicName]; ok {
		topic.RemoveClientChannel(clientID)
	}
	dis.RUnlock()
}

func (dis *Dispatcher) NotifyNodeMessage(topicNames []string, action string, body interface{}) {
	rpcBody, err := json.Marshal(body)
	if err != nil {
		return
	}

	dis.RLock()
	defer dis.RUnlock()

	for server := range dis.nodeServerMap {
		con, err := grpc.Dial(server, grpc.WithInsecure())
		if err != nil {
			continue
		}

		cli := protos.NewPublishClient(con)
		_, _ = cli.Publish(context.Background(), &protos.PublishMessage{
			Topics: topicNames,
			Action: action,
			Body:   rpcBody,
		})
	}
}
