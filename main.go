package main

import (
	"flag"
	"fmt"
	"gitlab.orayer.com/golang/pubsub/app"
	"os"
	"runtime"
)

var (
	// AppVersion 应用版本
	AppVersion string
	// BuildDate 构建日期
	BuildDate string
	// GitCommit 最后提交的git commit
	GitCommit string

	c string
	h bool
	v bool
)

func init() {
	flag.StringVar(&c, "c", "./config.toml", "config file path")
	flag.BoolVar(&h, "h", false, "this help")
	flag.BoolVar(&v, "v", false, "this help")

	flag.Parse()

	if h {
		help()
	}

	if v {
		version()
	}
}

func help() {
	fmt.Printf(`Real-time publish and subscribe service

Usage:
  pubsub [options]

Flags:
  -h,	help for gitlab.orayer.com/golang/pubsub
  -v,   print version
`)
	os.Exit(0)
}

func version() {
	fmt.Printf(`
Version: %s
GO Version: %s
Commit: %s
BuildTime: %s
`, AppVersion, runtime.Version(), GitCommit, BuildDate)
	os.Exit(0)
}

func main() {
	app.New(c).Run()
}
