package core

import (
	"fmt"
	"github.com/elliotcourant/noahdb/pkg/drivers/rqliter"
	"github.com/elliotcourant/noahdb/pkg/frunk"
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
	GetTable(name string) (Table, bool, error)
	GetTables(...string) ([]Table, error)
	GetColumns(tableId uint64) ([]Column, error)
	GetPrimaryKeyColumnByName(name string) (Column, bool, error)
	GetShardColumn(uint64) (Column, error)
	GetTablesInSchema(schema string, names ...string) ([]Table, error)
	GetTenantTable() (Table, bool, error)
}

func (ctx *base) Tables() TableContext {
	return &tableContext{
		ctx,
	}
}

func (ctx *tableContext) GetTenantTable() (Table, bool, error) {
	compiledSql, _, _ := getTablesQuery.
		Where(goqu.Ex{
			"table_type": TableType_Tenant,
		}).
		Limit(1).
		ToSql()

	rows, err := ctx.db.Query(compiledSql)
	if err != nil {
		return Table{}, false, err
	}

	result, err := ctx.tablesFromRows(rows)
	if err != nil {
		return Table{}, false, err
	}
	// No tenants table exists.
	if len(result) == 0 {
		return Table{}, false, nil
	}
	return result[0], true, nil
}

func (ctx *tableContext) GetTable(name string) (Table, bool, error) {
	tables, err := ctx.GetTables(name)
	if err != nil {
		return Table{}, false, err
	}
	if len(tables) == 0 {
		return Table{}, false, nil
	}
	return tables[0], true, nil
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

func (ctx *tableContext) GetColumns(tableId uint64) ([]Column, error) {
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

func (ctx *tableContext) GetPrimaryKeyColumnByName(tableName string) (Column, bool, error) {
	compiledSql, _, _ := getColumnsQuery.
		InnerJoin(
			goqu.I("tables"),
			goqu.On(goqu.I("tables.table_id").Eq(goqu.I("columns.table_id"))),
		).
		Where(goqu.Ex{
			"tables.table_name":   tableName,
			"columns.primary_key": true,
		}).
		Limit(1).
		ToSql()
	rows, err := ctx.db.Query(compiledSql)
	if err != nil {
		return Column{}, false, err
	}
	columns, err := ctx.columnsFromRows(rows)
	if err != nil {
		return Column{}, false, err
	}
	if len(columns) == 0 {
		return Column{}, false, nil
	}
	return columns[0], true, nil
}

func (ctx *tableContext) GetShardColumn(tableId uint64) (Column, error) {
	compiledSql, _, _ := getColumnsQuery.
		Where(goqu.Ex{
			"table_id":  tableId,
			"shard_key": true,
		}).Limit(1).ToSql()
	rows, err := ctx.db.Query(compiledSql)
	if err != nil {
		return Column{}, err
	}
	columns, err := ctx.columnsFromRows(rows)
	if len(columns) != 1 {
		return Column{}, fmt.Errorf("tried to find one shard column, found %d", len(columns))
	}
	return columns[0], err
}

func (ctx *tableContext) tablesFromRows(response *frunk.QueryResponse) ([]Table, error) {
	rows := rqliter.NewRqlRows(response)
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

func (ctx *tableContext) columnsFromRows(response *frunk.QueryResponse) ([]Column, error) {
	rows := rqliter.NewRqlRows(response)
	items := make([]Column, 0)
	for rows.Next() {
		item := Column{}
		if err := rows.Scan(
			&item.ColumnID,
			&item.TableID,
			&item.Type,
			&item.Sort,
			&item.ColumnName,
			&item.PrimaryKey,
			&item.Nullable,
			&item.ShardKey,
			&item.Serial,
			&item.ForeignColumnID); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
