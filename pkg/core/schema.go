package core

import (
	"fmt"
	"github.com/readystock/goqu"
	"strings"
)

type schemaContext struct {
	*base
}

type SchemaContext interface {
	GetSchemas() ([]Schema, error)
	Exists(string) (bool, error)
	NewSchema(string) (Schema, error)
}

func (ctx *base) Schema() SchemaContext {
	return &schemaContext{
		ctx,
	}
}

func (ctx *schemaContext) GetSchemas() ([]Schema, error) {
	return nil, nil
}

func (ctx *schemaContext) Exists(name string) (bool, error) {
	name = ctx.cleanSchemaName(name)

	sql, _, _ := goqu.
		From("schemas").
		Select(goqu.COUNT("schema_id")).
		Where(
			goqu.Ex{"schema_name": name}).
		ToSql()
	response, err := ctx.db.Query(sql)
	if err != nil {
		return false, err
	}
	return exists(response)
}

// NewSchema creates a new schema in noahdb.
func (ctx *schemaContext) NewSchema(name string) (schema Schema, err error) {
	name = ctx.cleanSchemaName(name)
	if ok, err := ctx.Exists(name); err != nil {
		return schema, fmt.Errorf("could not verify if there was a conflicting schema: %s", err.Error())
	} else if ok {
		return schema, fmt.Errorf("a schema with name [%s] already exists", name)
	}

	id, err := ctx.db.NextSequenceValueById(schemaIdSequencePath)
	if err != nil {
		return schema, err
	}

	sql := goqu.From("schemas").
		Insert(goqu.Record{
			"schema_id":   id,
			"schema_name": name,
		}).Sql
	_, err = ctx.db.Exec(sql)
	schema.SchemaID = id
	schema.SchemaName = name
	return schema, err
}

func (ctx *schemaContext) cleanSchemaName(name string) string {
	return strings.TrimSpace(strings.ToLower(name))
}
