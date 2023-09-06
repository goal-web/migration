package tests

import (
	"github.com/goal-web/console"
	"github.com/goal-web/contracts"
)

func NewService() contracts.ServiceProvider {
	return console.NewService(NewKernel)
}

func NewKernel(app contracts.Application) contracts.Console {
	return &Kernel{Kernel: console.NewKernel(app, []contracts.CommandProvider{}), app: app}
}

type Kernel struct {
	*console.Kernel
	app contracts.Application
}

func (kernel *Kernel) Schedule(schedule contracts.Schedule) {

}
