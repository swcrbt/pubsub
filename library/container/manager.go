package container

import (
	"gitlab.orayer.com/golang/pubsub/config"
	"gitlab.orayer.com/golang/pubsub/library/storage"
	"gitlab.orayer.com/golang/pubsub/library/dispatcher"
	"log"
	"os"
	"time"
)

var Mgr *Manager

type Manager struct {
	Config *config.Config
	Logger *log.Logger
	Storager *storage.Redis
	Dispatcher *dispatcher.Dispatcher
}

func NewManager(configFile string) (*Manager, error) {
	conf := config.LoadConfig(configFile)

	Mgr = &Manager{
		Config: conf,
		Storager: storage.NewRedis(conf.Storage.Address, conf.Storage.Password),
		Dispatcher: dispatcher.New(),
	}

	if conf.Logger.Type == "file" {
		fileName := conf.Logger.Target + time.Now().Format("20060102") + ".log"
		logIo, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
		if err != nil {
			return nil, err
		}

		Mgr.Logger = log.New(logIo, "", log.Ldate|log.Ltime|log.Lshortfile)
	} else {
		Mgr.Logger = log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lshortfile)
	}

	return Mgr, nil
}