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

func NewMigrate(app contracts.Application) contracts.Command {
	dir, _ := os.Getwd()
	if str, exists := app.Get("migrations.dir").(string); exists && str != "" {
		dir += "/" + str
	} else {
		dir += "/migrations"
	}
	return &Migrate{
		Command: commands.Base("migrate", "execute migrations"),
		conn:    app.Get("db").(contracts.DBConnection),
		dir:     dir,
	}
}

type Migrate struct {
	commands.Command
	conn contracts.DBConnection
	dir  string
}

type MigrateMsg struct {
	Batch  int           `json:"batch"`
	Path   string        `json:"path"`
	Action string        `json:"action"`
	Time   time.Duration `json:"time"`
}

func (cmd Migrate) Handle() any {
	logs.Default().Info("执行迁移")
	initTable(cmd.conn)

	var batch int
	if Migrations().Count() > 0 {
		batch = int(Migrations().Max("batch"))
	}

	var dir = cmd.StringOptional("path", cmd.dir)
	var items []MigrateMsg
	var migrated = Migrations().Get()
	var files = collection.New(getFiles(dir)).Filter(func(i int, s string) bool {
		return !strings.HasSuffix(s, ".down.sql") && migrated.Filter(func(i int, m *Migration) bool {
			return m.Path == s
		}).Count() == 0
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
			Batch:  batch + 1,
			Path:   path,
			Time:   time.Now().Sub(now),
		})
		Migrations().Create(contracts.Fields{
			"batch":      batch + 1,
			"path":       path,
			"created_at": carbon.Now().ToDateTimeString(),
		})
	}

	table.Output(items)

	return nil
}
