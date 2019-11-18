package app

import (
	"go-issued-service/service"
	"go-issued-service/library/container"
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
	app.Use(service.NewIssued())

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
}

func (app *App) Use(ser service.Service)  {
	app.services = append(app.services, ser)
}