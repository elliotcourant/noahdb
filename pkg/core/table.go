package core

import (
	"database/sql"
	"github.com/readystock/goqu"
)

var (
	getTablesQuery = goqu.
			From("tables").
			Select("tables.*")

	getColumnsQuery = goqu.
			From("columns").
			Select("columns.*")
)

type tableContext struct {
	*base
}

type TableContext interface {
	GetTables(...string) ([]Table, error)
}

func (ctx *base) Tables() TableContext {
	return &tableContext{
		ctx,
	}
}

func (ctx *tableContext) GetTables(names ...string) ([]Table, error) {
	query := getTablesQuery
	if len(names) > 0 {
		query = query.Where(goqu.Ex{
			"table_name": names,
		})
	}
	compiledSql, _, _ := query.ToSql()
	rows, err := ctx.db.Query(compiledSql)
	if err != nil {
		return nil, err
	}
	return ctx.tablesFromRows(rows)
}

func (ctx *tableContext) GetTablesInSchema(schema string, names ...string) ([]Table, error) {
	query := getTablesQuery.
		InnerJoin(goqu.I("schemas"), goqu.On(goqu.I("schemas.schema_id").Eq(goqu.I("tables.schema_id")))).
		Where(goqu.Ex{
			"schemas.schema_name": schema,
		})
	if len(names) > 0 {
		query = query.Where(goqu.Ex{
			"table_name": names,
		})
	}
	compiledSql, _, _ := query.ToSql()
	rows, err := ctx.db.Query(compiledSql)
	if err != nil {
		return nil, err
	}
	return ctx.tablesFromRows(rows)
}

func (ctx *tableContext) GetTable(schemaId int, tableName string) (Table, error) {
	return Table{}, nil
}

func (ctx *tableContext) GetColumns(tableId int) ([]Column, error) {
	compileSql, _, _ := getColumnsQuery.
		Where(goqu.Ex{
			"table_id": tableId,
		}).ToSql()
	rows, err := ctx.db.Query(compileSql)
	if err != nil {
		return nil, err
	}
	return ctx.columnsFromRows(rows)
}

func (ctx *tableContext) tablesFromRows(rows *sql.Rows) ([]Table, error) {
	defer rows.Close()
	items := make([]Table, 0)
	for rows.Next() {
		item := Table{}
		if err := rows.Scan(
			&item.TableID,
			&item.SchemaID,
			&item.TableName,
			&item.TableType); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (ctx *tableContext) columnsFromRows(rows *sql.Rows) ([]Column, error) {
	defer rows.Close()
	items := make([]Column, 0)
	for rows.Next() {
		item := Column{}
		if err := rows.Scan(
			&item.ColumnID,
			&item.TableID,
			&item.TypeID,
			&item.Sort,
			&item.ColumnName,
			&item.PrimaryKey,
			&item.Nullable,
			&item.ShardKey,
			&item.Serial); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
