package service

import (
	"gitlab.orayer.com/golang/pubsub/library/container"
	"gitlab.orayer.com/golang/pubsub/protos"
	"google.golang.org/grpc"
	"net"
)

type RpcService struct {
	server  *grpc.Server
}

func NewRpcService() *RpcService {
	return &RpcService{}
}

func (rser *RpcService) Run() error {
	lis, err := net.Listen("tcp", container.Mgr.Config.Server.RpcService.Address)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()

	protos.RegisterNodeServer(grpcServer, &NodeService{})
	protos.RegisterPublishServer(grpcServer, &PublishService{})

	go func() {
		container.Mgr.Logger.Printf("\"%s\" Server Run At: \"%s\"\n", rser.GetName(), container.Mgr.Config.Server.RpcService.Address)

		if err := grpcServer.Serve(lis); err != nil {
			container.Mgr.Logger.Printf("\"%s\" Server error: %v\n", rser.GetName(), err)
		}
	}()

	rser.server = grpcServer

	return nil
}

func (rser *RpcService) GetName() string {
	return "rpc-service"
}

func (rser *RpcService) Stop() error {
	if rser.server != nil {
		rser.server.GracefulStop()
	}
	return nil
}
