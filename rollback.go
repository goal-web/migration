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

func NewRollback() (contracts.Command, contracts.CommandHandlerProvider) {
	return commands.Base("migrate:rollback", "execute models.MigrationQuery"),
		func(app contracts.Application) contracts.CommandHandler {
			return &Rollback{
				conn: app.Get("db").(contracts.DBConnection),
				dir:  getDir(app.Get("config").(contracts.Config)),
			}
		}
}

type Rollback struct {
	commands.Command
	conn contracts.DBConnection
	dir  string
}

func (cmd Rollback) Handle() any {
	logs.Default().Info("执行回滚")
	initTable(cmd.conn)

	var batch int
	if models.MigrationQuery().Count() > 0 {
		batch = cmd.IntOptional("batch", int(models.MigrationQuery().Max("batch")))
	} else {
		return nil
	}
	var items []MigrateMsg
	var migrated = models.MigrationQuery().Where("batch", batch).Get()
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
			Batch:  batch,
			Path:   migration.Path,
			Action: "rollback",
			Time:   time.Since(now),
		})
		models.MigrationQuery().Where("id", migration.Id).Delete()
	})

	table.Output(items)

	return nil
}
