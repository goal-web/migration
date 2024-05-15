package migration

import (
	"github.com/goal-web/database/table"
	"github.com/goal-web/supports/class"
)

var Class = class.Make[Migration]()

func Migrations() *table.Table[Migration] {
	return table.Class(Class, "migrations")
}

type Migration struct {
	table.Model[Migration] `json:"-"`

	Id        string `json:"id"`
	Path      string `json:"path"`
	Batch     int    `json:"batch"`
	CreatedAt string `json:"created_at"`
}
