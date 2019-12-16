package service

import (
	"gitlab.orayer.com/golang/pubsub/library/container"
	"gitlab.orayer.com/golang/pubsub/protos"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net"
)

type PublishService struct {
}

type RpcPublisher struct {
	handler *PublishService
	server  *grpc.Server
}

func NewRpcPublisher() *RpcPublisher {
	return &RpcPublisher{
		handler: &PublishService{},
	}
}

func (rec *RpcPublisher) Run() error {
	lis, err := net.Listen("tcp", container.Mgr.Config.Server.PublisherRpc.Address)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()

	protos.RegisterPublishServer(grpcServer, rec.handler)

	go func() {
		container.Mgr.Logger.Printf("\"%s\" Server Run At: \"%s\"\n", rec.GetName(), container.Mgr.Config.Server.PublisherRpc.Address)

		if err := grpcServer.Serve(lis); err != nil {
			container.Mgr.Logger.Printf("\"%s\" Server error: %v\n", rec.GetName(), err)
		}
	}()

	rec.server = grpcServer

	return nil
}

func (rec *RpcPublisher) GetName() string {
	return "publisher-rpc"
}

func (rec *RpcPublisher) Stop() error {
	if rec.server != nil {
		rec.server.GracefulStop()
	}
	return nil
}

func (ser *PublishService) Release(ctx context.Context, req *protos.PublishBody) (*protos.PublishResponse, error) {
	var resp []*protos.PublishResponse_ResponseBody

	/*result := container.Mgr.Dispatcher.Publish(req.Topics, req.Action, req.Body)

	result.Range(func(topic, res interface{}) bool {

		res.(*sync.Map).Range(func(cid, isReply interface{}) bool {
			resp = append(resp, &protos.PublishResponse_ResponseBody{Topic: topic.(string), ChannelId: cid.(string), IsReply: isReply.(bool)})
			return true
		})

		return true
	})*/

	container.Mgr.Dispatcher.Publish(req.Topics, req.Action, req.Body)
	return &protos.PublishResponse{Body: resp}, nil
}
