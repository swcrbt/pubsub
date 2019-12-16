package service_test

import (
	"bytes"
	"encoding/json"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.orayer.com/golang/pubsub/library/container"
	"gitlab.orayer.com/golang/pubsub/service"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"testing"
	"time"
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

func BenchmarkHttpPublisher(b *testing.B) {
	_, err := container.NewManager("../config.toml")
	if err != nil {
		container.Mgr.Logger.Printf("register container fail: %v \n", err)
		return
	}

	rand.Seed(time.Now().UnixNano())
	max := 30
	idPrefix := "zxc_"
	var wg sync.WaitGroup

	for i:=0; i<max; i++ {
		id := strconv.Itoa(i)
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			randTopic := idPrefix + strconv.Itoa(int(rand.Int31n(int32(200))))
			container.Mgr.Logger.Printf("\"%s\" publish topic: %s", id, randTopic)
			reqBody := service.PublishBody{
				Topics: []string{"xxx", randTopic},
				Action: randTopic,
				Body:   map[string]string{"a": "b", "b": "c"},
			}

			data, err := json.Marshal(reqBody)
			if err != nil {
				container.Mgr.Logger.Printf("\"%s\" json marshal fail: %v \n", id, err)
				return
			}

			req := bytes.NewBuffer(data)
			u := url.URL{
				Scheme: "http",
				//Host: "192.168.30.227:9992",
				Host: "127.0.0.1:8888",
				Path: "publish",
			}

			resp, err := http.Post(u.String(), "application/json", req)
			if err != nil {
				container.Mgr.Logger.Printf("\"%s\" http request fail: %v \n", id, err)
				return
			}
			defer resp.Body.Close()

			container.Mgr.Logger.Printf("\"%s\" publish resp statuscode: %v data: %v", id,  resp.StatusCode, resp.Body)
		}(id)
	}

	wg.Wait()
}
