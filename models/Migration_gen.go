// Code generated by goal-pro. DO NOT EDIT.
// versions:
//
//	goal-pro v1.0.0
//	go       v1.23
//
// updated_at: 2024-10-01 15:15:02
// source: migration.proto
package models

import (
	"github.com/goal-web/contracts"
	"github.com/goal-web/database/table"
	"github.com/goal-web/supports/utils"
)

func NewMigrationModel(fields contracts.Fields) *MigrationModel {
	var model MigrationModel
	model.Set(fields)
	return &model
}

func MigrationQuery() *table.Table[MigrationModel] {
	return table.NewQuery("migrations", NewMigrationModel).SetPrimaryKey("id")
}

type MigrationModel struct {
	Id int32 `json:"id"`

	Path string `json:"path"`

	Batch int32 `json:"batch"`

	CreatedAt string `json:"created_at"`

	_update contracts.Fields
}

func (model *MigrationModel) Exists() bool {
	return MigrationQuery().Where("id", model.GetPrimaryKey()).Count() > 0
}

func (model *MigrationModel) Save() contracts.Exception {
	if model._update == nil {
		return nil
	}
	if MigrationModelSaving != nil {
		if err := MigrationModelSaving(model); err != nil {
			return err
		}
	}
	_, err := MigrationQuery().Where("id", model.GetPrimaryKey()).UpdateE(model._update)
	if err == nil {
		model._update = nil
		if MigrationModelSaved != nil {
			MigrationModelSaved(model)
		}
	}

	return err
}

func (model *MigrationModel) Set(fields contracts.Fields) {
	for key, value := range fields {
		if key == "id" {
			switch v := value.(type) {
			case int32:
				model.SetId(v)
			case func() int32:
				model.SetId(v())
			}
		}
		if key == "path" {
			switch v := value.(type) {
			case string:
				model.SetPath(v)
			case func() string:
				model.SetPath(v())
			}
		}
		if key == "batch" {
			switch v := value.(type) {
			case int32:
				model.SetBatch(v)
			case func() int32:
				model.SetBatch(v())
			}
		}
		if key == "created_at" {
			switch v := value.(type) {
			case string:
				model.SetCreatedAt(v)
			case func() string:
				model.SetCreatedAt(v())
			}
		}
	}
}

func (model *MigrationModel) Only(key ...string) contracts.Fields {
	var fields = make(contracts.Fields)
	for _, k := range key {
		if k == "id" {
			fields[k] = model.GetId()
			continue
		}
		if k == "path" {
			fields[k] = model.GetPath()
			continue
		}
		if k == "batch" {
			fields[k] = model.GetBatch()
			continue
		}
		if k == "created_at" {
			fields[k] = model.GetCreatedAt()
			continue
		}

		if MigrationModelAppends[k] != nil {
			fields[k] = MigrationModelAppends[k](model)
		}
	}
	return fields
}

func (model *MigrationModel) Except(keys ...string) contracts.Fields {
	var excepts = map[string]struct{}{}
	for _, k := range keys {
		excepts[k] = struct{}{}
	}
	var fields = make(contracts.Fields)
	for key, value := range model.ToFields() {
		if _, ok := excepts[key]; ok {
			continue
		}
		fields[key] = value
	}
	return fields
}

var MigrationModelAppends = map[string]func(model *MigrationModel) any{}

func (model *MigrationModel) ToFields() contracts.Fields {
	fields := contracts.Fields{
		"id":         model.GetId(),
		"path":       model.GetPath(),
		"batch":      model.GetBatch(),
		"created_at": model.GetCreatedAt(),
	}

	for key, f := range MigrationModelAppends {
		fields[key] = f(model)
	}

	return fields
}

func (model *MigrationModel) Update(fields contracts.Fields) contracts.Exception {

	if MigrationModelUpdating != nil {
		if err := MigrationModelUpdating(model, fields); err != nil {
			return err
		}
	}

	if model._update != nil {
		utils.MergeFields(model._update, fields)
	}

	_, err := MigrationQuery().Where("id", model.GetPrimaryKey()).UpdateE(fields)

	if err == nil {
		model.Set(fields)
		model._update = nil
		if MigrationModelUpdated != nil {
			MigrationModelUpdated(model, fields)
		}
	}

	return err
}

func (model *MigrationModel) Refresh() contracts.Exception {
	fields, err := table.ArrayQuery("migrations").Where("id", model.GetPrimaryKey()).FirstE()
	if err != nil {
		return err
	}

	model.Set(*fields)
	return nil
}

func (model *MigrationModel) Delete() contracts.Exception {

	if MigrationModelDeleting != nil {
		if err := MigrationModelDeleting(model); err != nil {
			return err
		}
	}

	_, err := MigrationQuery().Where("id", model.GetPrimaryKey()).DeleteE()
	if err == nil && MigrationModelDeleted != nil {
		MigrationModelDeleted(model)
	}

	return err
}

var (
	MigrationModelIdGetter         func(model *MigrationModel, raw int32) int32
	MigrationModelIdSetter         func(model *MigrationModel, raw int32) int32
	MigrationModelPathGetter       func(model *MigrationModel, raw string) string
	MigrationModelPathSetter       func(model *MigrationModel, raw string) string
	MigrationModelBatchGetter      func(model *MigrationModel, raw int32) int32
	MigrationModelBatchSetter      func(model *MigrationModel, raw int32) int32
	MigrationModelCreatedAtGetter  func(model *MigrationModel, raw string) string
	MigrationModelCreatedAtSetter  func(model *MigrationModel, raw string) string
	MigrationModelSaving           func(model *MigrationModel) contracts.Exception
	MigrationModelSaved            func(model *MigrationModel)
	MigrationModelUpdating         func(model *MigrationModel, fields contracts.Fields) contracts.Exception
	MigrationModelUpdated          func(model *MigrationModel, fields contracts.Fields)
	MigrationModelDeleting         func(model *MigrationModel) contracts.Exception
	MigrationModelDeleted          func(model *MigrationModel)
	MigrationModelPrimaryKeyGetter func(model *MigrationModel) any
)

func (model *MigrationModel) GetPrimaryKey() any {
	if MigrationModelPrimaryKeyGetter != nil {
		return MigrationModelPrimaryKeyGetter(model)
	}

	return model.Id
}

func (model *MigrationModel) GetId() int32 {
	if MigrationModelIdGetter != nil {
		return MigrationModelIdGetter(model, model.Id)
	}
	return model.Id
}

func (model *MigrationModel) SetId(value int32) {
	if MigrationModelIdSetter != nil {
		value = MigrationModelIdSetter(model, value)
	}

	if model._update == nil {
		model._update = contracts.Fields{"id": value}
	} else {
		model._update["id"] = value
	}
	model.Id = value
}

func (model *MigrationModel) GetPath() string {
	if MigrationModelPathGetter != nil {
		return MigrationModelPathGetter(model, model.Path)
	}
	return model.Path
}

func (model *MigrationModel) SetPath(value string) {
	if MigrationModelPathSetter != nil {
		value = MigrationModelPathSetter(model, value)
	}

	if model._update == nil {
		model._update = contracts.Fields{"path": value}
	} else {
		model._update["path"] = value
	}
	model.Path = value
}

func (model *MigrationModel) GetBatch() int32 {
	if MigrationModelBatchGetter != nil {
		return MigrationModelBatchGetter(model, model.Batch)
	}
	return model.Batch
}

func (model *MigrationModel) SetBatch(value int32) {
	if MigrationModelBatchSetter != nil {
		value = MigrationModelBatchSetter(model, value)
	}

	if model._update == nil {
		model._update = contracts.Fields{"batch": value}
	} else {
		model._update["batch"] = value
	}
	model.Batch = value
}

func (model *MigrationModel) GetCreatedAt() string {
	if MigrationModelCreatedAtGetter != nil {
		return MigrationModelCreatedAtGetter(model, model.CreatedAt)
	}
	return model.CreatedAt
}

func (model *MigrationModel) SetCreatedAt(value string) {
	if MigrationModelCreatedAtSetter != nil {
		value = MigrationModelCreatedAtSetter(model, value)
	}

	if model._update == nil {
		model._update = contracts.Fields{"created_at": value}
	} else {
		model._update["created_at"] = value
	}
	model.CreatedAt = value
}
