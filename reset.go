package migration

import (
	"fmt"
	"github.com/goal-web/contracts"
	"github.com/goal-web/migration/models"
	"github.com/goal-web/supports/commands"
	"github.com/goal-web/supports/logs"
	"github.com/modood/table"
	"os"
	"strings"
	"time"
)

func NewReset() (contracts.Command, contracts.CommandHandlerProvider) {
	return commands.Base("migrate:reset", "Rollback all database migrations"), func(app contracts.Application) contracts.CommandHandler {
		return &Reset{
			conn: app.Get("db").(contracts.DBConnection),
			dir:  getDir(app.Get("config").(contracts.Config)),
		}
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
	var migrated = models.MigrationQuery().Get()
	var dir = cmd.StringOptional("path", cmd.dir)

	migrated.Map(func(i int, migration *models.MigrationModel) {
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
			Batch:  int(migration.Batch),
			Path:   migration.Path,
			Action: "reset",
			Time:   time.Since(now),
		})
		models.MigrationQuery().Where("id", migration.Id).Delete()
	})

	table.Output(items)

	return nil
}
