package service

import (
	"gitlab.orayer.com/golang/issue/library/container"
	"gitlab.orayer.com/golang/issue/protos"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net"
)

type ReleaseService struct {
}

type RpcReceiver struct {
	handler *ReleaseService
	server *grpc.Server
}

func NewRpcReceiver() *RpcReceiver {
	return &RpcReceiver{
		handler: &ReleaseService{},
	}
}

func (rec *RpcReceiver) Run() error {
	lis, err := net.Listen("tcp", container.Mgr.Config.Server.ReceiverRpc.Address)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()

	protos.RegisterIReleaseServiceServer(grpcServer, rec.handler)

	go func() {
		container.Mgr.Logger.Printf("\"%s\" Server Run At: \"%s\"\n", rec.GetName(), container.Mgr.Config.Server.ReceiverRpc.Address)

		if err := grpcServer.Serve(lis); err != nil {
			container.Mgr.Logger.Printf("\"%s\" Server error: %v\n", rec.GetName(), err)
		}
	}()

	rec.server = grpcServer

	return nil
}

func (rec *RpcReceiver) GetName() string {
	return "receiver-rpc"
}

func (rec *RpcReceiver) Stop() error {
	if rec.server != nil {
		rec.server.GracefulStop()
	}
	return nil
}

func (ser *ReleaseService) Release(ctx context.Context, req *protos.ReleaseBody) (*protos.ReleaseResponse, error) {
	container.Mgr.Dispatcher.Push(req.Action, req.UniqIds, req.Data)

	return &protos.ReleaseResponse{Value: map[string]string{}}, nil
}