package service_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go-issued-service/library/container"
	"testing"
)

var _ = BeforeSuite(func() {
	_, err := container.NewManager("../config.toml")
	Expect(err).NotTo(HaveOccurred())
})

func TestService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Issue Service")
}