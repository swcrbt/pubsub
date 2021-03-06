package service_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.orayer.com/golang/pubsub/library/container"
	"testing"
)

var _ = BeforeSuite(func() {
	_, err := container.NewManager("../config.toml")
	Expect(err).NotTo(HaveOccurred())
})

func TestService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "subscribe Service")
}