package tests

import (
	"fmt"
	"github.com/goal-web/console/inputs"
	"github.com/goal-web/contracts"
	"regexp"
	"testing"
)

func TestMigrate(t *testing.T) {
	app := initApp()

	app.Call(func(console contracts.Console) {
		console.Run(inputs.StringArrayInput{ArgsArray: []string{"migrate", "--path=migrations"}})
	})
}

func TestRollback(t *testing.T) {
	app := initApp()

	app.Call(func(console contracts.Console) {
		console.Run(inputs.StringArrayInput{ArgsArray: []string{"migrate:rollback", "--path=migrations"}})
	})
}

func TestRefresh(t *testing.T) {
	app := initApp()

	app.Call(func(console contracts.Console) {
		console.Run(inputs.StringArrayInput{ArgsArray: []string{"migrate:refresh", "--path=migrations"}})
	})
}

func TestStatus(t *testing.T) {
	app := initApp()

	app.Call(func(console contracts.Console) {
		console.Run(inputs.StringArrayInput{ArgsArray: []string{"migrate:status"}})
	})
}

func TestReset(t *testing.T) {
	app := initApp()

	app.Call(func(console contracts.Console) {
		console.Run(inputs.StringArrayInput{ArgsArray: []string{"migrate:reset"}})
	})
}

func TestMakeMigration(t *testing.T) {
	app := initApp()

	app.Call(func(console contracts.Console) {
		console.Run(inputs.StringArrayInput{ArgsArray: []string{"make:migration", "create_posts_table"}})
	})
}

func TestRegexp(t *testing.T) {
	reg, _ := regexp.Compile("create_(.*)_table")
	fmt.Println(reg.FindStringSubmatch("create_posts_table")[1])
}
