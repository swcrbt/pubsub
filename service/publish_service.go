package service

import (
	"context"
	"encoding/json"
	"gitlab.orayer.com/golang/errors"
	"gitlab.orayer.com/golang/pubsub/library/container"
	"gitlab.orayer.com/golang/pubsub/protos"
	"time"
)

type PublishService struct {
}

func (ps *PublishService) Publish (ctx context.Context, msg *protos.PublishMessage) (*protos.PublishResponse, error) {
	var body map[string]interface{}
	err := json.Unmarshal(msg.Body, &body)
	if err != nil {
		return nil, errors.New("publisher/params_error", "invalid request payload")
	}
	container.Mgr.Logger.Printf("[RPC-PUB] %v | %s | %v | %v | %v \n",
		time.Now().Format("2006/01/02 - 15:04:05"),
		container.Mgr.Config.Server.RpcService.Address,
		msg.Topics,
		msg.Action,
		body,
	)
	container.Mgr.Dispatcher.Publish(msg.Topics, msg.Action, body)
	return &protos.PublishResponse{}, nil
}