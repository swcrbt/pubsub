package app

import (
	"github.com/gin-gonic/gin"
	"gitlab.orayer.com/golang/pubsub/library/container"
	"gitlab.orayer.com/golang/pubsub/service"
	"os"
	"os/signal"
	"runtime"
	"time"
)

type App struct {
	services []service.Service
}

func New (configFile string) *App {
	if _ , err := container.NewManager(configFile); err != nil {
		panic(err)
	}

	container.Mgr.RegisterNode(container.Mgr.Config.Server.RpcService.Address)

	gin.SetMode(container.Mgr.Config.Server.Mode)

	app := &App{}
	app.Use(service.NewSubscriber())
	app.Use(service.NewPublisher())
	app.Use(service.NewRpcService())

	//if gin.IsDebugging() {
	//	app.Use(service.NewPProf())
	//}

	return app
}

func (app *App) Run() {
	for _, ser := range app.services {
		if err := ser.Run(); err != nil {
			container.Mgr.Logger.Printf("\"%s\" run failed: %v\n", ser.GetName(), err)
		}
	}

	exitChan := make(chan byte)

	go func() {
		var rtm runtime.MemStats
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				runtime.ReadMemStats(&rtm)

				container.Mgr.Logger.Printf(`
numgc: %v, pausetotalns: %v, numgoroutine: %d, cpunum: %d
alloc: %v, totalalloc: %v, sys: %v, mallocs: %v, frees: %v, liveobjects: %v
`,
rtm.NumGC, rtm.PauseTotalNs, runtime.NumGoroutine(), runtime.NumCPU(),
rtm.Alloc, rtm.TotalAlloc, rtm.Sys, rtm.Mallocs, rtm.Frees, (rtm.Mallocs - rtm.Frees),
				)
			case <-exitChan:
				return
			}
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	close(exitChan)

	container.Mgr.Logger.Printf("quit signal")
	container.Mgr.UnRegisterNode(container.Mgr.Config.Server.RpcService.Address)
	app.Stop()
}

func (app *App) Reload()  {
	for _, ser := range app.services {
		if err := ser.Stop(); err != nil {
			container.Mgr.Logger.Printf("\"%s\" run failed: %v\n", ser.GetName(), err)
		}
		if err := ser.Run(); err != nil {
			container.Mgr.Logger.Printf("\"%s\" stop failed: %v\n", ser.GetName(), err)
		}
	}
}

func (app *App) Stop ()  {
	for _, ser := range app.services {
		if err := ser.Stop(); err != nil {
			container.Mgr.Logger.Printf("\"%s\" stop failed: %v\n", ser.GetName(), err)
		} else {
			container.Mgr.Logger.Printf("\"%s\" stop successed\n", ser.GetName())
		}
	}
}

func (app *App) Use(ser service.Service)  {
	app.services = append(app.services, ser)
}