package service

import (
	//_ "github.com/mkevac/debugcharts"
	"gitlab.orayer.com/golang/issue/library/container"
	"net/http"
	_ "net/http/pprof"
)

type PProf struct {
}

func NewPProf() *PProf {
	return &PProf{}
}

func (p *PProf) Run () error {
	go func() {
		container.Mgr.Logger.Printf("\"%s\" Server Run At: \"%s\"\n", p.GetName(), container.Mgr.Config.Server.PProf.Address)

		container.Mgr.Logger.Println(http.ListenAndServe(container.Mgr.Config.Server.PProf.Address, nil))
	}()

	return nil
}

func (p *PProf) GetName() string {
	return "pprof"
}

func (p *PProf) Stop() error {
	return nil
}