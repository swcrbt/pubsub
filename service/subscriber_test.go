package service_test

import (
	"encoding/base64"
	"encoding/json"
	"github.com/gorilla/websocket"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.orayer.com/golang/pubsub/library/container"
	"gitlab.orayer.com/golang/pubsub/middleware"
	"net/url"
	"strconv"
	"sync"
	"testing"
	"time"
)

var _ = BeforeSuite(func() {
	_, err := container.NewManager("../config.toml")
	Expect(err).NotTo(HaveOccurred())
})

func TestService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "subscribe Service")
}

func BenchmarkService(b *testing.B) {
	_, err := container.NewManager("../config.toml")
	if err != nil {
		b.Logf("register container fail: %v \n", err)
		return
	}

	max := 200
	keyPrefix := "sub_auth_key_"
	secret := "test"
	var wg sync.WaitGroup

	for i:=0; i<max; i++ {
		id := strconv.Itoa(i)
		val := middleware.StorageRecord{
			CryptoType: middleware.CryptoTypeTicket,
			ChannelID: id,
			Topics: []string{"xxx", "zxc_" + id},
			Secret: secret,
		}
		if data, err := json.Marshal(val); err == nil {
			container.Mgr.Logger.Println(keyPrefix + id, b.N)
			if err = container.Mgr.Storager.Set(keyPrefix + id, data); err != nil {
				container.Mgr.Logger.Printf("set storage fail: %v \n", err)
				continue
			}

			key, err := json.Marshal(middleware.AuthRecord{Key:id, Secret:secret})
			if err != nil {
				container.Mgr.Logger.Printf("json auth fail: %v \n", err)
				continue
			}

			wg.Add(1)
			go func(id string, key []byte) {
				defer wg.Done()
				u := url.URL{
					Scheme: "ws",
					//Host: "192.168.30.227:9999",
					Host: "127.0.0.1:9999",
					Path: "subscribe",
					RawQuery: "key=" + base64.StdEncoding.EncodeToString(key),
				}

				c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
				if err != nil {
					container.Mgr.Logger.Printf("\"%s\" webscoker connected fail: %v \n", id, err)
					return
				}

				defer c.Close()
				done := make(chan struct{})

				go func() {
					defer close(done)
					for {
						msgType, message, err := c.ReadMessage()
						if err != nil {
							container.Mgr.Logger.Printf("\"%s\" read err: %v \n", id, err)
							return
						}
						container.Mgr.Logger.Printf("\"%s\" read msg type: %v, message: %v \n", id, msgType, string(message))
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
							container.Mgr.Logger.Printf("\"%s\" write fail: %v \n", id, err)
							return
						}
					}
				}
			}(id, key)
		}
	}

	wg.Wait()
}