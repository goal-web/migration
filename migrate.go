package migration

import (
	"fmt"
	"github.com/goal-web/collection"
	"github.com/goal-web/contracts"
	"github.com/goal-web/migration/models"
	"github.com/goal-web/supports/commands"
	"github.com/goal-web/supports/logs"
	"github.com/golang-module/carbon/v2"
	"github.com/modood/table"
	"os"
	"strings"
	"time"
)

func NewMigrate() (contracts.Command, contracts.CommandHandlerProvider) {
	return commands.Base("migrate", "execute migrations"),
		func(app contracts.Application) contracts.CommandHandler {
			return &Migrate{
				conn: app.Get("db").(contracts.DBConnection),
				dir:  getDir(app.Get("config").(contracts.Config)),
			}
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
	if models.MigrationQuery().Count() > 0 {
		batch = int(models.MigrationQuery().Max("batch"))
	}

	var dir = cmd.StringOptional("path", cmd.dir)
	var items []MigrateMsg
	var migrated = models.MigrationQuery().Get()
	var files = collection.New(getFiles(dir)).Filter(func(i int, s string) bool {
		return !strings.HasSuffix(s, ".down.sql") && migrated.Filter(func(i int, m *models.MigrationModel) bool {
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
			Time:   time.Since(now),
		})
		models.MigrationQuery().Create(contracts.Fields{
			"batch":      batch + 1,
			"path":       path,
			"created_at": carbon.Now().ToDateTimeString(),
		})
	}

	table.Output(items)

	return nil
}
