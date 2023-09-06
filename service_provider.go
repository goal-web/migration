package migration

import (
	"github.com/goal-web/contracts"
)

type serviceProvider struct {
}

func NewService() contracts.ServiceProvider {
	return &serviceProvider{}
}

func (s serviceProvider) Register(app contracts.Application) {
	app.Call(func(console contracts.Console) {
		for _, provider := range []contracts.CommandProvider{
			NewMigrate,
			NewRollback,
			NewRefresh,
			NewGenerator,
			NewShowStatus,
			NewReset,
		} {
			console.RegisterCommand(provider(app).GetName(), provider)
		}
	})
}

func (s serviceProvider) Start() error {
	return nil
}

func (s serviceProvider) Stop() {
}
