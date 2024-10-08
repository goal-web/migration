package migrate

import (
	"database/sql"
	"fmt"
	"github.com/goal-web/collection"
	"github.com/goal-web/contracts"
	"github.com/goal-web/supports/logs"
	"github.com/goal-web/supports/utils"
	"reflect"
	"strings"
)

type Migrator func(executor contracts.SqlExecutor) contracts.Exception

// ColumnInfo 结构体用于存储表的列信息
type ColumnInfo struct {
	Field   string         `db:"Field" json:"Field"`
	Type    string         `db:"Type" json:"Type"`
	Null    string         `db:"Null" json:"Null"`
	Key     string         `db:"Key" json:"Key"`
	Default sql.NullString `db:"Default" json:"Default"`
	Extra   string         `db:"Extra" json:"Extra"`
}

// IndexInfo 用于接收 MySQL 表索引信息的结构体
type IndexInfo struct {
	Table        string `db:"Table"`        // 表名
	NonUnique    any    `db:"Non_unique"`   // 是否唯一（0 表示唯一，1 表示非唯一）
	KeyName      string `db:"Key_name"`     // 索引名称
	SeqInIndex   int    `db:"Seq_in_index"` // 索引中列的顺序（1 表示第一个列）
	ColumnName   string `db:"Column_name"`
	Collation    any    `db:"Collation"`
	Cardinality  any    `db:"Cardinality"`
	SubPart      any    `db:"Sub_part"`
	Packed       any    `db:"Packed"`
	Null         any    `db:"Null"`
	IndexType    any    `db:"Index_type"`
	Comment      any    `db:"Comment"`
	IndexComment any    `db:"Index_comment"`
	Visible      any    `db:"Visible"`
	Expression   any    `db:"Expression"`
}

func Auto(factory contracts.DBFactory, migrators ...Migrator) {
	logs.Default().Info("开始执行自动迁移")
	exception := factory.Connection().Transaction(func(executor contracts.SqlExecutor) contracts.Exception {
		var exception contracts.Exception
		for _, migrator := range migrators {
			if exception = migrator(executor); exception != nil {
				return exception
			}
		}

		return exception
	})

	if exception != nil {
		logs.Default().Error("自动迁移失败")
		logs.Default().Error(exception.Error())
	} else {
		logs.Default().Info("已完成所有自动迁移")
	}
}

func Migrate(tableName string, indexes []string, model any, executor contracts.SqlExecutor) contracts.Exception {
	var fields []ColumnInfo
	var statements []string

	exception := executor.Select(&fields, fmt.Sprintf("describe `%s`", tableName))
	if exception != nil {
		if !strings.HasSuffix(exception.Error(), " doesn't exist") {
			return exception
		}
		createStatement := fmt.Sprintf("create table if not exists `%s`", tableName)
		var tableFields []string

		utils.EachStructField(reflect.ValueOf(model), model, func(field reflect.StructField, value reflect.Value) {
			db := field.Tag.Get("db")
			if field.IsExported() && db != "" {
				var fieldName, dbType, constraints string
				for i, tag := range strings.Split(db, ";") {
					if i == 0 {
						fieldName = tag
					} else if strings.HasPrefix(tag, "type:") {
						dbType = strings.TrimPrefix(tag, "type:")
					} else {
						constraints += " " + tag
					}
				}
				tableFields = append(tableFields, fmt.Sprintf("`%s` %s %s;", fieldName, dbType, constraints))
			}
		})

		createStatement += fmt.Sprintf(" (%s)", strings.Join(tableFields, ", "))
		statements = append(statements, createStatement)
	}

	if len(fields) > 0 {
		fieldsGrouped := collection.New(fields).Pluck("field")
		utils.EachStructField(reflect.ValueOf(model), model, func(field reflect.StructField, value reflect.Value) {
			db := field.Tag.Get("db")
			if field.IsExported() && db != "" {
				if strings.Contains(strings.ToLower(db), "primary key") {
					return
				}
				var fieldName, dbType string
				var constraints []string
				for i, tag := range strings.Split(db, ";") {
					if i == 0 {
						fieldName = tag
					} else if strings.HasPrefix(tag, "type:") {
						dbType = strings.ToLower(strings.TrimPrefix(tag, "type:"))
					} else {
						constraints = append(constraints, strings.ToLower(tag))
					}
				}

				constraintsStr := strings.Join(constraints, " ")
				if existsField, exists := fieldsGrouped[fieldName]; exists {
					var changed bool
					for _, constraint := range constraints {
						if (constraint == "not null" && existsField.Null == "YES") || (constraint == "null" && existsField.Null == "NO") {
							changed = true
						}
					}
					if strings.ToLower(existsField.Type) != dbType {
						changed = true
					}
					for _, constraint := range constraints {
						if strings.HasPrefix(constraint, "default ") {
							existsDefault := strings.ToLower(existsField.Default.String)
							newDefault := strings.TrimSuffix(strings.TrimPrefix(constraint, "default "), " on update current_timestamp")
							if existsDefault != newDefault && fmt.Sprintf("'%s'", existsDefault) != newDefault {
								changed = true
							}
						}
					}

					if changed {
						statements = append(statements, fmt.Sprintf("ALTER TABLE %s MODIFY COLUMN %s %s %s;", tableName, fieldName, dbType, constraintsStr))
					}
				} else {
					statements = append(statements, fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s %s;", tableName, fieldName, dbType, constraintsStr))
				}
			}
		})
	}

	var existsIndexes []IndexInfo
	exception = executor.Select(&existsIndexes, fmt.Sprintf("show index from %s", tableName))
	if exception != nil {
		return exception
	}

	indexesGrouped := collection.New(existsIndexes).GroupBy("key_name")
	for _, index := range indexes {
		indexData := strings.Split(index, ";")
		if _, exists := indexesGrouped[indexData[1]]; !exists {
			statements = append(statements, fmt.Sprintf("create %s %s on %s %s;", indexData[0], indexData[1], tableName, strings.ReplaceAll(indexData[2], ";", ",")))
		}
	}

	if len(statements) > 0 {
		query := strings.Join(statements, "\n")
		_, exception = executor.Exec(query)
		if exception != nil {
			return exception
		}
		logs.Default().Info(fmt.Sprintf("%s 已完成迁移.", tableName))
		logs.Default().Info(fmt.Sprintf("%s 迁移内容：%s", tableName, query))
	} else {
		logs.Default().Info(fmt.Sprintf("%s 无需执行迁移.", tableName))
	}

	return nil
}
