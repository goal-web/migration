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

func NewRollback(app contracts.Application) contracts.Command {
	return &Rollback{
		Command: commands.Base("migrate:rollback", "execute migrations"),
		conn:    app.Get("db").(contracts.DBConnection),
		dir:     getDir(app.Get("config").(contracts.Config)),
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
	if Migrations().Count() > 0 {
		batch = cmd.IntOptional("batch", int(Migrations().Max("batch")))
	} else {
		return nil
	}
	var items []MigrateMsg
	var migrated = Migrations().Where("batch", batch).Get()
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
			Batch:  batch,
			Path:   migration.Path,
			Action: "rollback",
			Time:   time.Now().Sub(now),
		})
		Migrations().Where("id", migration.Id).Delete()
	})

	table.Output(items)

	return nil
}
