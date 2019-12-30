package container

import (
	"context"
	"encoding/json"
	"gitlab.orayer.com/golang/pubsub/config"
	"gitlab.orayer.com/golang/pubsub/dispatcher"
	"gitlab.orayer.com/golang/pubsub/library/storage"
	"gitlab.orayer.com/golang/pubsub/protos"
	"google.golang.org/grpc"
	"log"
	"os"
	"time"
)

const NODE_SERVER_KEY = "subpub_node_servers"

var Mgr *Manager

type Manager struct {
	Config *config.Config
	Logger *log.Logger
	Storager *storage.Redis
	Dispatcher *dispatcher.Dispatcher
}

func NewManager(configFile string) (*Manager, error) {
	conf := config.LoadConfig(configFile)

	storager := storage.NewRedis(conf.Storage.Address, conf.Storage.Password)

	Mgr = &Manager{
		Config: conf,
		Storager: storager,
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

	Mgr.UpdateNodeServer()

	return Mgr, nil
}

func (mgr *Manager) UpdateNodeServer()  {
	nodeServer := make(map[string]*dispatcher.Node)

	if nodeData, _ := mgr.Storager.Get(NODE_SERVER_KEY); len(nodeData) > 0 {
		_ = json.Unmarshal(nodeData, &nodeServer)
	}

	delete(nodeServer, mgr.Config.Server.RpcService.Address)

	mgr.Dispatcher.SetNodeServer(nodeServer)
}

func (mgr *Manager) RegisterNode (server string) {
	nodeServer := make(map[string]*dispatcher.Node)

	nodeData, _ := mgr.Storager.Get(NODE_SERVER_KEY)
	if len(nodeData) > 0 {
		_ = json.Unmarshal(nodeData, &nodeServer)
	}

	nodeServer[server] = dispatcher.NewNode(server)

	if data, err := json.Marshal(nodeServer); err == nil {
		_ = mgr.Storager.Set(NODE_SERVER_KEY, data)
	}

	for ser := range nodeServer {
		if ser == server {
			continue
		}
		_ = mgr.NotifyNodeRefresh(ser)
	}
}

func (mgr *Manager) UnRegisterNode (server string) {
	nodeServer := make(map[string]*dispatcher.Node)

	nodeData, _ := mgr.Storager.Get(NODE_SERVER_KEY)
	if len(nodeData) > 0 {
		_ = json.Unmarshal(nodeData, &nodeServer)
	}

	delete(nodeServer, server)

	mgr.Dispatcher.SetNodeServer(nodeServer)

	if data, err := json.Marshal(nodeServer); err == nil {
		_ = mgr.Storager.Set(NODE_SERVER_KEY, data)
	}

	for ser := range nodeServer {
		_ = mgr.NotifyNodeRefresh(ser)
	}
}

func (mgr *Manager) NotifyNodeRefresh (server string) error {
	con, err := grpc.Dial(server, grpc.WithInsecure())
	if err != nil {
		return err
	}
	
	cli := protos.NewNodeClient(con)
	_, err = cli.RefreshNode(context.Background(), &protos.NodeRequest{})
	if err != nil {
		return err
	}

	return nil
}