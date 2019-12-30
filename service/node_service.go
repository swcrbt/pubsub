package service

import (
	"context"
	"gitlab.orayer.com/golang/pubsub/library/container"
	"gitlab.orayer.com/golang/pubsub/protos"
	"time"
)

type NodeService struct {
}

func (ns *NodeService) RefreshNode (ctx context.Context, r *protos.NodeRequest) (*protos.NodeResponse, error) {
	container.Mgr.Logger.Printf("[RPC-RENODE] %v | %s \n",
		time.Now().Format("2006/01/02 - 15:04:05"),
		container.Mgr.Config.Server.RpcService.Address,
	)
	container.Mgr.UpdateNodeServer()
	return &protos.NodeResponse{}, nil
}