package service_test

import (
	"bytes"
	"encoding/json"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"go-issued-service/library/container"
	"go-issued-service/service"
	"net/http"
	"strconv"
)

var _ = Describe("receiver-http", func() {
	It("Test http release", func() {
		httpServer := service.NewHttpReceiver()
		err := httpServer.Run();
		Expect(err).NotTo(HaveOccurred())

		Expect(httpServer.GetName()).To(Equal("receiver-http"))

		body := service.ReleaseBody{
			Action:  "pgygame",
			UniqIds: []string{"1", "2", "3"},
			Data:    map[string]string{"a": "b", "b": "c"},
		}

		data, err := json.Marshal(body)
		Expect(err).NotTo(HaveOccurred())

		req := bytes.NewBuffer(data)

		resp, err := http.Post("http://127.0.0.1:"+strconv.Itoa(container.Mgr.Config.Server.ReceiverHttp.Port)+"/release",
			"application/json", req)
		Expect(err).NotTo(HaveOccurred())

		Expect(resp.StatusCode).To(Equal(http.StatusOK))
	})
})
