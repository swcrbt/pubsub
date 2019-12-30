package config

import (
	"io/ioutil"
	"time"

	"github.com/BurntSushi/toml"
)

// ServerConfig server config
type ServerConfig struct {
	Mode            string
	ShutdownTimeout time.Duration `toml:"shutdown_timeout"`
	Subscriber      SubscriberConfig
	Publisher       PublisherConfig
	RpcService      RpcServiceConfig `toml:"rpc_service"`
	PProf           PProfConfig
}

type SubscriberConfig struct {
	Address string

	ReadDeadline  time.Duration `toml:"read_deadline"`
	WriteDeadline time.Duration `toml:"write_deadline"`
}

type PublisherConfig struct {
	Address string
}

type RpcServiceConfig struct {
	Address string
}

type PProfConfig struct {
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
