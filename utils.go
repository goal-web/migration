package migration

import (
	"github.com/goal-web/contracts"
	"os"
	"path/filepath"
	"strings"
)

func initTable(connection contracts.DBConnection) {
	_, e := connection.Exec(Table)
	if e != nil {
		panic(e)
	}
}

func getFiles(dir string) []string {
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

const Table = "CREATE TABLE IF NOT EXISTS migrations\n(\n    `id`       INT UNSIGNED AUTO_INCREMENT,\n    path       varchar(255),\n    batch      int,\n    created_at timestamp,\n    PRIMARY KEY (`id`)\n) ENGINE = InnoDB\n  DEFAULT CHARSET = utf8mb4;"
