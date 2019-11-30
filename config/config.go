package config

import (
	"io/ioutil"
	"time"

	"github.com/BurntSushi/toml"
)

// ServerConfig server config
type ServerConfig struct {
	Mode          string
	Subscriber    SubscriberConfig
	PublisherHttp PublisherHttpConfig `toml:"publisher_http"`
	PublisherRpc  PublisherRpcConfig  `toml:"publisher_rpc"`
	PProf         PProfConfig
}

type SubscriberConfig struct {
	Port int

	ReadDeadline      time.Duration `toml:"read_deadline"`
	WriteDeadline     time.Duration `toml:"write_deadline"`
}

type PublisherHttpConfig struct {
	Port int
}

type PublisherRpcConfig struct {
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
