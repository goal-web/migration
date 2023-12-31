package migration

import (
	"fmt"
	"github.com/goal-web/collection"
	"github.com/goal-web/contracts"
	"github.com/goal-web/supports/commands"
	"github.com/goal-web/supports/logs"
	"github.com/modood/table"
	"os"
	"strings"
	"time"
)

func NewRefresh(app contracts.Application) contracts.Command {
	dir, _ := os.Getwd()
	if str, exists := app.Get("migrations.dir").(string); exists && str != "" {
		dir += "/" + str
	} else {
		dir += "/migrations"
	}
	return &Refresh{
		Command: commands.Base("migrate:refresh", "execute migrations"),
		conn:    app.Get("db").(contracts.DBConnection),
		dir:     dir,
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
		Migrations().Get().Map(func(i int, migration Migration) {
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
			Time:   time.Now().Sub(now),
		})
		Migrations().Create(contracts.Fields{
			"batch":      1,
			"path":       path,
			"created_at": time.Now(),
		})
	}

	table.Output(items)

	return nil
}
