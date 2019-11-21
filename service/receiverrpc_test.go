package service_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go-issued-service/library/container"
	"go-issued-service/protos"
	"go-issued-service/service"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var _ = Describe("receiver-rpc", func() {
	It("Test rpc release", func() {
		rpcServer := service.NewRpcReceiver()
		err := rpcServer.Run();
		Expect(err).NotTo(HaveOccurred())

		Expect(rpcServer.GetName()).To(Equal("receiver-rpc"))

		con, err := grpc.Dial(container.Mgr.Config.Server.ReceiverRpc.Address, grpc.WithInsecure())
		Expect(err).NotTo(HaveOccurred())

		cli := protos.NewIReleaseServiceClient(con)

		_, err = cli.Release(
			context.Background(),
			&protos.ReleaseBody{
				Action:"pgygame",
				UniqIds:[]string{"1","2","3"},
				Data: map[string]string{"a": "b", "b": "c"},
			})

		Expect(err).NotTo(HaveOccurred())
	})
})