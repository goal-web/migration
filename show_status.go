package migration

import (
	"github.com/goal-web/contracts"
	"github.com/goal-web/supports/commands"
	"github.com/goal-web/supports/logs"
	"github.com/modood/table"
)

func NewShowStatus(app contracts.Application) contracts.Command {
	return &ShowStatus{
		Command: commands.Base("migrate:reset", "Rollback all database migrations"),
		conn:    app.Get("db").(contracts.DBConnection),
	}
}

type ShowStatus struct {
	commands.Command
	conn contracts.DBConnection
	dir  string
}

func (cmd ShowStatus) init() {
	_, e := cmd.conn.Exec(Table)
	if e != nil {
		panic(e)
	}
}

func (cmd ShowStatus) Handle() any {
	cmd.init()
	items := Migrations().OrderByDesc("batch").OrderByDesc("id").Get().ToArray()

	if len(items) > 0 {
		table.Output(items)
	} else {
		logs.Default().Info("删除")
	}

	return nil
}
