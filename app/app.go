package app

import (
	"github.com/gin-gonic/gin"
	"gitlab.orayer.com/golang/issue/library/container"
	"gitlab.orayer.com/golang/issue/service"
	"os"
	"os/signal"
)

type App struct {
	services []service.Service
}

func New (configFile string) *App {
	if _ , err := container.NewManager(configFile); err != nil {
		panic(err)
	}

	app := &App{}
	app.Use(service.NewIssuer())
	app.Use(service.NewHttpReceiver())
	//app.Use(service.NewRpcReceiver())

	if gin.IsDebugging() {
		app.Use(service.NewPProf())
	}

	return app
}

func (app *App) Run() {
	for _, ser := range app.services {
		if err := ser.Run(); err != nil {
			container.Mgr.Logger.Printf("\"%s\" run failed: %v\n", ser.GetName(), err)
		}
	}

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

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