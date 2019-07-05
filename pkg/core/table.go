package core

import (
	"fmt"
	"github.com/elliotcourant/noahdb/pkg/drivers/rqliter"
	"github.com/elliotcourant/noahdb/pkg/frunk"
	"github.com/readystock/golog"
	"github.com/readystock/goqu"
	"strings"
	"time"
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
	NewTable(table Table, columns []Column) (Table, []Column, error)
	NextSequenceID(table Table, column Column) (uint64, error)
	GetTable(name string) (Table, bool, error)
	GetTables(...string) ([]Table, error)
	GetColumns(tableId uint64) ([]Column, error)
	GetPrimaryKeyColumnByName(name string) (Column, bool, error)
	GetSequenceColumnForTable(tableId uint64) (Column, bool, error)
	GetShardColumn(uint64) (Column, error)
	GetTablesInSchema(schema string, names ...string) ([]Table, error)
	GetTenantTable() (Table, bool, error)
}

func (ctx *base) Tables() TableContext {
	return &tableContext{
		ctx,
	}
}

func (ctx *tableContext) NewTable(table Table, columns []Column) (Table, []Column, error) {
	// exists, err := ctx.Exists(table.TableName)
	tableId, err := ctx.db.NextSequenceValueById(tableIdSequencePath)
	if err != nil {
		return Table{}, nil, err
	}
	table.TableID = tableId

	for i := range columns {
		colId, err := ctx.db.NextSequenceValueById(columnIdSequencePath)
		if err != nil {
			return Table{}, nil, err
		}
		columns[i].TableID, columns[i].ColumnID = tableId, colId
	}

	insertTableSql := goqu.
		From("tables").
		Insert(goqu.Record{
			"table_id":     table.TableID,
			"schema_id":    1,
			"table_name":   table.TableName,
			"table_type":   table.TableType,
			"has_sequence": table.HasSequence,
		}).Sql

	columnSql := make([]string, len(columns))
	for i, col := range columns {
		fcid := &col.ForeignColumnID
		if *fcid == 0 {
			fcid = nil
		}
		columnSql[i] = goqu.
			From("columns").
			Insert(goqu.Record{
				"column_id":         col.ColumnID,
				"table_id":          col.TableID,
				"type_id":           col.Type,
				"sort":              col.Sort,
				"column_name":       col.ColumnName,
				"primary_key":       col.PrimaryKey,
				"nullable":          col.Nullable,
				"shard_key":         col.ShardKey,
				"serial":            col.Serial,
				"foreign_column_id": fcid,
			}).Sql
	}

	compiledSql := strings.Join(append([]string{insertTableSql}, columnSql...), ";\n")

	_, err = ctx.db.Exec(compiledSql)
	return table, columns, err
}

func (ctx *tableContext) NextSequenceID(table Table, column Column) (uint64, error) {
	startTimestamp := time.Now()
	defer func() {
		golog.Verbosef("[%s] new sequence id for [%s.%s]", time.Since(startTimestamp), table.TableName, column.ColumnName)
	}()
	return ctx.db.NextSequenceValueById(fmt.Sprintf("/schema/%d/table/%d/column/%d/sequence", table.SchemaID, column.TableID, column.ColumnID))
}

func (ctx *tableContext) Exists(name string) (bool, error) {
	compiledSql, _, _ := goqu.
		From("tables").
		Select(
			goqu.COUNT("tables.table_id"),
		).
		Where(goqu.Ex{
			"tables.table_name": name,
		}).
		ToSql()
	rows, err := ctx.db.Query(compiledSql)
	if err != nil {
		return false, err
	}
	return exists(rows)
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

func (ctx *tableContext) GetSequenceColumnForTable(tableId uint64) (Column, bool, error) {
	compiledSql, _, _ := getColumnsQuery.
		Where(goqu.Ex{
			"table_id": tableId,
			"serial":   true,
		}).Limit(1).ToSql()
	rows, err := ctx.db.Query(compiledSql)
	if err != nil {
		return Column{}, false, err
	}
	columns, err := ctx.columnsFromRows(rows)
	if len(columns) > 1 {
		return Column{}, false, fmt.Errorf("tried to find one shard column, found %d", len(columns))
	}
	if len(columns) == 0 {
		return Column{}, false, nil
	}
	return columns[0], true, err
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
			&item.TableType,
			&item.HasSequence); err != nil {
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
		var foreignColumnId *uint64
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
			&foreignColumnId); err != nil {
			return nil, err
		}
		if foreignColumnId != nil {
			item.ForeignColumnID = *foreignColumnId
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
