package config_test

import (
	"testing"

	"gitlab.orayer.com/golang/pubsub/config"
)

func Test_LoadConfig(t *testing.T) {
	cfile := "../config.toml"

	c := config.LoadConfig(cfile)

	t.Logf("%#v", c.Server.Mode)
}
