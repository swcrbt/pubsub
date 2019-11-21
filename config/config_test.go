package config_test

import (
	"testing"

	"go-issued-service/config"
)

func Test_LoadConfig(t *testing.T) {
	cfile := "../config.toml"

	c := config.LoadConfig(cfile)

	t.Logf("%#v", c.Server.Mode)
}
