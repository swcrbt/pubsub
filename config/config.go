package config

import (
	"io/ioutil"
	"time"

	"github.com/BurntSushi/toml"
)

// ServerConfig server config
type ServerConfig struct {
	Mode         string
	Issuer       IssuerConfig
	ReceiverHttp ReceiverHttpConfig `toml:"receiver_http"`
	ReceiverRpc  ReceiverRpcConfig  `toml:"receiver_rpc"`
}

type IssuerConfig struct {
	Port int

	HeartbeatInterval time.Duration `toml:"heartbeat_interval"`
	HeartbeatTimeout  int           `toml:"heartbeat_timeout"`
}

type ReceiverHttpConfig struct {
	Port int
}

type ReceiverRpcConfig struct {
	Address string
}

type StorageConfig struct {
	Address  string
	Password string
}

type LoggerConfig struct {
	Type   string
	Target string
}

// Config app config
type Config struct {
	Server  ServerConfig
	Storage StorageConfig
	Logger  LoggerConfig
}

// load config from toml file
func LoadConfig(filename string) *Config {

	c, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	config := &Config{}

	if err := toml.Unmarshal(c, config); err != nil {
		panic(err)
	}

	return config
}
