package config_test

import (
	"testing"

	"storage.oray.com/config"
)

func Test_MustLoadConfig(t *testing.T) {

	cfile := "../config.toml"

	config := config.MustLoadConfig(cfile)

	t.Logf("%#v", config.Upload["profile"])
}
