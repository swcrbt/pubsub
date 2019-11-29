package service_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.orayer.com/golang/issue/library/container"
	"gitlab.orayer.com/golang/issue/protos"
	"gitlab.orayer.com/golang/issue/service"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var _ = Describe("Publisher-rpc", func() {
	It("Test rpc release", func() {
		rpcServer := service.NewRpcPublisher()
		err := rpcServer.Run();
		Expect(err).NotTo(HaveOccurred())

		Expect(rpcServer.GetName()).To(Equal("publisher-rpc"))

		con, err := grpc.Dial(container.Mgr.Config.Server.PublisherRpc.Address, grpc.WithInsecure())
		Expect(err).NotTo(HaveOccurred())

		cli := protos.NewPublishClient(con)

		_, err = cli.Release(
			context.Background(),
			&protos.PublishBody{
				Topics: []string{"pgygame_173060"},
				Action: "xxx",
				Body:   map[string]string{"a": "b", "b": "c"},
			})

		Expect(err).NotTo(HaveOccurred())
	})
})