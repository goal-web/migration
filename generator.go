package migration

import (
	"fmt"
	"github.com/goal-web/contracts"
	"github.com/goal-web/supports/commands"
	"os"
	"regexp"
	"time"
)

func NewGenerator(app contracts.Application) contracts.Command {
	dir, _ := os.Getwd()
	if str, exists := app.Get("migrations.dir").(string); exists && str != "" {
		dir += "/" + str
	} else {
		dir += "/migrations"
	}
	return &Generator{
		Command: commands.Base("make:migration {name}", "Create a new migration file"),
		dir:     dir,
	}
}

type Generator struct {
	commands.Command
	dir string
}

func (cmd Generator) Handle() any {
	var name = time.Now().Format("2006_01_02_150405") + "_" + cmd.GetString("name")
	var dir = cmd.StringOptional("path", cmd.dir)
	reg, _ := regexp.Compile("create_(.*)")
	upSql := ""
	downSql := ""

	if results := reg.FindStringSubmatch(name); len(results) > 0 {
		tableName := results[1]
		upSql = fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s"+
			"("+
			"    `id`       INT UNSIGNED AUTO_INCREMENT,"+
			"    created_at timestamp,"+
			"    updated_at timestamp,"+
			"    PRIMARY KEY (`id`)"+
			"    ) ENGINE = InnoDB"+
			"    DEFAULT CHARSET = utf8mb4;", tableName)

		downSql = "drop table if exists " + tableName
	}

	err := os.WriteFile(fmt.Sprintf("%s/%s.sql", dir, name), []byte(upSql), os.ModePerm)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(fmt.Sprintf("%s/%s.down.sql", dir, name), []byte(downSql), os.ModePerm)
	if err != nil {
		panic(err)
	}

	return nil
}
