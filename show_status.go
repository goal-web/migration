package migration

import (
	"github.com/goal-web/contracts"
	"github.com/goal-web/migration/models"
	"github.com/goal-web/supports/commands"
	"github.com/modood/table"
	"strings"
)

func NewShowStatus() (contracts.Command, contracts.CommandHandlerProvider) {
	return commands.Base("migrate:status", "Rollback all database models.MigrationQuery"), func(app contracts.Application) contracts.CommandHandler {
		return &ShowStatus{
			conn: app.Get("db").(contracts.DBConnection),
			dir:  getDir(app.Get("config").(contracts.Config)),
		}
	}
}

type ShowStatus struct {
	commands.Command
	conn contracts.DBConnection
	dir  string
}

type Status struct {
	Path   string `json:"path"`
	Batch  int    `json:"batch"`
	Status string `json:"status"`
}

func (cmd ShowStatus) Handle() any {
	initTable(cmd.conn)

	items := models.MigrationQuery().OrderByDesc("batch").OrderByDesc("id").Get().Pluck("path")

	var dir = cmd.StringOptional("path", cmd.dir)
	var list []Status

	for _, path := range getFiles(dir) {
		item, exists := items[path]
		text := "pending"
		if exists {
			text = "migrated"
		}
		if !strings.HasSuffix(path, ".down.sql") {
			list = append(list, Status{
				Path:   path,
				Batch:  int(item.Batch),
				Status: text,
			})
		}
	}
	table.Output(list)

	return nil
}
