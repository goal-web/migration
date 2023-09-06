package tests

import (
	"github.com/goal-web/console"
	"github.com/goal-web/contracts"
	"github.com/goal-web/migration"
)

func NewService() contracts.ServiceProvider {
	return console.NewService(NewKernel)
}

func NewKernel(app contracts.Application) contracts.Console {
	return &Kernel{Kernel: console.NewKernel(app, []contracts.CommandProvider{
		migration.NewMigrate,
		migration.NewRollback,
		migration.NewRefresh,
		migration.NewGenerator,
		migration.NewShowStatus,
		migration.NewReset,
	}), app: app}
}

type Kernel struct {
	*console.Kernel
	app contracts.Application
}

func (kernel *Kernel) Schedule(schedule contracts.Schedule) {

}
