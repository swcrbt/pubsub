package service_test

import (
	"bytes"
	"encoding/json"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.orayer.com/golang/pubsub/library/container"
	"gitlab.orayer.com/golang/pubsub/service"
	"net/http"
)

var _ = Describe("Publisher", func() {
	It("Test http release", func() {
		httpServer := service.NewPublisher()
		err := httpServer.Run();
		Expect(err).NotTo(HaveOccurred())

		Expect(httpServer.GetName()).To(Equal("publisher"))

		body := service.PublishMessage{
			Topics: []string{"pgygame_173060"},
			Action: "xxx",
			Body:   map[string]string{"a": "b", "b": "c"},
		}

		data, err := json.Marshal(body)
		Expect(err).NotTo(HaveOccurred())

		req := bytes.NewBuffer(data)

		resp, err := http.Post("http://"+ container.Mgr.Config.Server.Publisher.Address +"/publish",
			"application/json", req)
		Expect(err).NotTo(HaveOccurred())

		Expect(resp.StatusCode).To(Equal(http.StatusOK))
	})
})
