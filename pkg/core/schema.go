package core

import (
	"github.com/readystock/goqu"
	"strings"
)

type schemaContext struct {
	*base
}

type SchemaContext interface {
	GetSchemas() ([]Schema, error)
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
	name = strings.ToLower(name)
	sql, _, _ := goqu.
		From("schemas").
		Select(goqu.COUNT("schema_id")).
		Where(
			goqu.Ex{"schema_name": name}).
		ToSql()
	count, err := ctx.db.Count(sql)
	if err != nil {
		return false, err
	}
	return count == 1, nil
}

// NewSchema creates a new schema in noahdb.
func (ctx *schemaContext) NewSchema(name string) (Schema, error) {
	return Schema{}, nil
}
