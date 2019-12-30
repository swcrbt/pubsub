package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"gitlab.orayer.com/golang/pubsub/service"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"sync"
)

var (
	num        = flag.Int("n", 1, "并发数")
	requestUrl = flag.String("u", "0.0.0.0:9992", "发布地址")
	topic      = flag.String("t", "test", "发布主题")
)

func main() {
	flag.Parse()

	var wg sync.WaitGroup

	for i := 0; i < *num; i++ {
		id := strconv.Itoa(i)
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			reqBody := service.PublishMessage{
				Topics: []string{*topic},
				Action: "test",
				Body: map[string]interface{}{
					"ts": "string",
					"ti": 123,
					"tb": true,
					"tl": []string{"tl1", "tl2"},
					"tm": map[string]string{"tmk": "tmv"},
				},
			}

			data, err := json.Marshal(reqBody)
			if err != nil {
				log.Printf("\"%s\" json marshal fail: %v \n", id, err)
				return
			}

			req := bytes.NewBuffer(data)
			u := url.URL{
				Scheme: "http",
				Host:   *requestUrl,
				Path:   "publish",
			}

			resp, err := http.Post(u.String(), "application/json", req)
			if err != nil {
				log.Printf("\"%s\" http request fail: %v \n", id, err)
				return
			}
			defer resp.Body.Close()

			log.Printf("\"%s\" publish resp statuscode: %v data: %v", id, resp.StatusCode, resp.Body)
		}(id)
	}

	wg.Wait()
}
