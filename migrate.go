package migration

import (
	"fmt"
	"github.com/goal-web/collection"
	"github.com/goal-web/contracts"
	"github.com/goal-web/supports/commands"
	"github.com/goal-web/supports/logs"
	"os"
	"path/filepath"
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

const Table = "CREATE TABLE IF NOT EXISTS migrations\n(\n    `id`       INT UNSIGNED AUTO_INCREMENT,\n    path       varchar(255),\n    batch      int,\n    created_at timestamp,\n    PRIMARY KEY (`id`)\n) ENGINE = InnoDB\n  DEFAULT CHARSET = utf8mb4;"

func (cmd Migrate) init() {
	_, e := cmd.conn.Exec(Table)
	if e != nil {
		panic(e)
	}
}

func (cmd Migrate) Files() []string {
	var dir = cmd.StringOptional("path", cmd.dir)
	var files []string
	fs, err := os.Stat(dir)
	if err != nil {
		panic(err)
	}

	if fs.IsDir() {
		err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if strings.HasSuffix(info.Name(), ".sql") {
				files = append(files, info.Name())
			}
			return nil
		})
		if err != nil {
			panic(err)
		}
	} else if strings.HasSuffix(dir, ".sql") {
		files = []string{fs.Name()}
	}

	return files
}

func (cmd Migrate) Handle() any {
	logs.Default().Info("执行迁移")
	cmd.init()

	var batch int
	if Migrations().Count() > 0 {
		batch = int(Migrations().Max("batch"))
	}

	var migrated = Migrations().Get()

	var files = collection.New(cmd.Files()).Filter(func(i int, s string) bool {
		return !strings.HasSuffix(s, "down.sql") && migrated.Where("path", s).Count() == 0
	}).ToArray()

	for _, path := range files {
		sqlBytes, err := os.ReadFile(fmt.Sprintf("%s/%s", cmd.dir, path))
		if err != nil {
			panic(err)
		}
		_, e := cmd.conn.Exec(string(sqlBytes))
		if e != nil {
			panic(e)
		}
		Migrations().Create(contracts.Fields{
			"batch":      batch + 1,
			"path":       path,
			"created_at": time.Now(),
		})
	}

	fmt.Println(files)

	return nil
}
