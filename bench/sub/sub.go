package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"github.com/gorilla/websocket"
	"gitlab.orayer.com/golang/pubsub/config"
	"gitlab.orayer.com/golang/pubsub/library/storage"
	"gitlab.orayer.com/golang/pubsub/middleware"
	"log"
	"net/url"
	"strconv"
	"sync"
	"time"
)

var (
	num        = flag.Int("n", 1, "并发数")
	requestUrl = flag.String("u", "0.0.0.0:9991", "订阅地址")
	topic      = flag.String("t", "test", "订阅主题")
	sid        = flag.String("id", "", "订阅客户端ID")
	configPath = flag.String("c", "./config.toml", "配置文件")
)

func main() {
	flag.Parse()

	conf := config.LoadConfig(*configPath)
	storager := storage.NewRedis(conf.Storage.Address, conf.Storage.Password)

	var wg sync.WaitGroup
	secret := "test"

	for i := 0; i < *num; i++ {
		id := strconv.Itoa(i)
		val := middleware.StorageRecord{
			CryptoType: middleware.CryptoTypeTicket,
			ID:         *sid,
			Topics:     []string{*topic},
			Secret:     secret,
		}
		if data, err := json.Marshal(val); err == nil {
			if err = storager.Set(middleware.AuthKeyPrefix+id, data); err != nil {
				log.Printf("set storage fail: %v \n", err)
				continue
			}

			key, err := json.Marshal(middleware.AuthRecord{Key: id, Secret: secret})
			if err != nil {
				log.Printf("json auth fail: %v \n", err)
				continue
			}

			wg.Add(1)
			go func(id string, key []byte) {
				defer wg.Done()
				u := url.URL{
					Scheme:   "ws",
					Host:     *requestUrl,
					Path:     "subscribe",
					RawQuery: "key=" + base64.StdEncoding.EncodeToString(key),
				}

				c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
				if err != nil {
					log.Printf("\"%s\" webscoker connected fail: %v \n", id, err)
					return
				}

				defer c.Close()
				done := make(chan struct{})

				go func() {
					defer close(done)
					for {
						msgType, message, err := c.ReadMessage()
						if err != nil {
							log.Printf("\"%s\" read err: %v \n", id, err)
							return
						}
						log.Printf("\"%s\" read msg type: %v, message: %v \n", id, msgType, string(message))
					}
				}()

				ticker := time.NewTicker(20 * time.Second)
				defer ticker.Stop()

				for {
					select {
					case <-done:
						return
					case <-ticker.C:
						err := c.WriteMessage(websocket.PingMessage, nil)
						if err != nil {
							log.Printf("\"%s\" write fail: %v \n", id, err)
							return
						}
					}
				}
			}(id, key)
		}
	}

	wg.Wait()
}
