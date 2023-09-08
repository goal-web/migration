package migration

import (
	"fmt"
	"github.com/goal-web/contracts"
	"github.com/goal-web/supports/commands"
	"github.com/goal-web/supports/logs"
	"github.com/modood/table"
	"os"
	"strings"
	"time"
)

func NewReset(app contracts.Application) contracts.Command {
	dir, _ := os.Getwd()
	if str, exists := app.Get("migrations.dir").(string); exists && str != "" {
		dir += "/" + str
	} else {
		dir += "/migrations"
	}
	return &Reset{
		Command: commands.Base("migrate:reset", "Rollback all database migrations"),
		conn:    app.Get("db").(contracts.DBConnection),
		dir:     dir,
	}
}

type Reset struct {
	commands.Command
	conn contracts.DBConnection
	dir  string
}

func (cmd Reset) Handle() any {
	logs.Default().Info("执行重置")
	initTable(cmd.conn)

	var items []MigrateMsg
	var migrated = Migrations().Get()
	var dir = cmd.StringOptional("path", cmd.dir)

	migrated.Map(func(i int, migration Migration) {
		sqlBytes, err := os.ReadFile(fmt.Sprintf("%s/%s", dir, strings.ReplaceAll(migration.Path, ".sql", ".down.sql")))
		if err != nil {
			panic(err)
		}
		now := time.Now()
		_, e := cmd.conn.Exec(string(sqlBytes))
		if e != nil {
			panic(e)
		}
		items = append(items, MigrateMsg{
			Batch:  migration.Batch,
			Path:   migration.Path,
			Action: "reset",
			Time:   time.Now().Sub(now),
		})
		Migrations().Where("id", migration.Id).Delete()
	})

	table.Output(items)

	return nil
}
