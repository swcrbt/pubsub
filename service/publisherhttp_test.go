package service_test

import (
	"bytes"
	"encoding/json"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.orayer.com/golang/issue/library/container"
	"gitlab.orayer.com/golang/issue/service"
	"net/http"
	"strconv"
)

var _ = Describe("Publisher-http", func() {
	It("Test http release", func() {
		httpServer := service.NewHttpPublisher()
		err := httpServer.Run();
		Expect(err).NotTo(HaveOccurred())

		Expect(httpServer.GetName()).To(Equal("publisher-http"))

		body := service.PublishBody{
			Topics: []string{"pgygame_173060"},
			Action: "xxx",
			Body:   map[string]string{"a": "b", "b": "c"},
		}

		data, err := json.Marshal(body)
		Expect(err).NotTo(HaveOccurred())

		req := bytes.NewBuffer(data)

		resp, err := http.Post("http://127.0.0.1:"+strconv.Itoa(container.Mgr.Config.Server.PublisherHttp.Port)+"/publish",
			"application/json", req)
		Expect(err).NotTo(HaveOccurred())

		Expect(resp.StatusCode).To(Equal(http.StatusOK))
	})
})
