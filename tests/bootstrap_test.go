package tests

import (
	"github.com/goal-web/console/inputs"
	"github.com/goal-web/contracts"
	"testing"
)

func TestBootstrap(t *testing.T) {
	app := initApp()

	app.Call(func(console contracts.Console) {
		console.Run(inputs.StringArrayInput{ArgsArray: []string{"migrate", "--path=migrations"}})
	})
}
