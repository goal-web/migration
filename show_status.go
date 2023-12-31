package migration

import (
	"github.com/goal-web/contracts"
	"github.com/goal-web/supports/commands"
	"github.com/modood/table"
	"os"
	"strings"
)

func NewShowStatus(app contracts.Application) contracts.Command {
	dir, _ := os.Getwd()
	if str, exists := app.Get("migrations.dir").(string); exists && str != "" {
		dir += "/" + str
	} else {
		dir += "/migrations"
	}
	return &ShowStatus{
		Command: commands.Base("migrate:status", "Rollback all database migrations"),
		conn:    app.Get("db").(contracts.DBConnection),
		dir:     dir,
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

	items := Migrations().OrderByDesc("batch").OrderByDesc("id").Get().Pluck("path")

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
				Batch:  item.Batch,
				Status: text,
			})
		}
	}
	table.Output(list)

	return nil
}
