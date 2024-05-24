package migration

import (
	"fmt"
	"github.com/goal-web/collection"
	"github.com/goal-web/contracts"
	"github.com/goal-web/supports/commands"
	"github.com/goal-web/supports/logs"
	"github.com/golang-module/carbon/v2"
	"github.com/modood/table"
	"os"
	"strings"
	"time"
)

func NewRefresh(app contracts.Application) contracts.Command {
	return &Refresh{
		Command: commands.Base("migrate:refresh", "execute migrations"),
		conn:    app.Get("db").(contracts.DBConnection),
		dir:     getDir(app.Get("config").(contracts.Config)),
	}
}

type Refresh struct {
	commands.Command
	conn contracts.DBConnection
	dir  string
}

func (cmd Refresh) Handle() any {
	logs.Default().Info("执行 refresh")
	initTable(cmd.conn)

	var items []MigrateMsg
	var dir = cmd.StringOptional("path", cmd.dir)
	if Migrations().Count() > 0 {
		var batch = cmd.IntOptional("batch", int(Migrations().Max("batch")))
		Migrations().Get().Map(func(i int, migration *Migration) {
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
			Migrations().Where("id", migration.Id).Delete()
		})
	}

	var files = collection.New(getFiles(dir)).Filter(func(i int, s string) bool {
		return !strings.HasSuffix(s, ".down.sql")
	}).ToArray()

	for _, path := range files {
		sqlBytes, err := os.ReadFile(fmt.Sprintf("%s/%s", dir, path))
		if err != nil {
			panic(err)
		}
		now := time.Now()
		_, e := cmd.conn.Exec(string(sqlBytes))
		if e != nil {
			panic(e)
		}
		items = append(items, MigrateMsg{
			Action: "migrate",
			Batch:  1,
			Path:   path,
			Time:   time.Since(now),
		})
		Migrations().Create(contracts.Fields{
			"batch":      1,
			"path":       path,
			"created_at": carbon.Now().ToDateTimeString(),
		})
	}

	table.Output(items)

	return nil
}
