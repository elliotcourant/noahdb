package core

import (
	"github.com/elliotcourant/noahdb/pkg/types"
	"github.com/elliotcourant/timber"
)

var (
	pgSchemas = []Schema{
		{
			SchemaName: "information_schema",
		},
		{
			SchemaName: "pg_catalog",
		},
	}
	pgTables = []struct {
		Table   Table
		Columns []Column
	}{
		{
			Table: Table{
				TableName: "pg_aggregate",
			},
			Columns: []Column{
				{
					ColumnName: "aggfnoid",
					Type:       types.Type_regproc,
				},
				{
					ColumnName: "aggkind",
					Type:       types.Type_char,
				},
				{
					ColumnName: "aggnumdirectargs",
					Type:       types.Type_int2,
				},
				{
					ColumnName: "aggtransfn",
					Type:       types.Type_regproc,
				},
				{
					ColumnName: "aggfinalfn",
					Type:       types.Type_regproc,
				},
			},
		},
	}
)

func (ctx *base) setupPostgresSystem() {
	timber.Debugf("creating %d postgres schema(s)", len(pgSchemas))
	for _, schema := range pgSchemas {
		s, err := ctx.Schema().NewSchema(schema.SchemaName)
		if err != nil {
			timber.Fatalf("failed to create schema [%s]: %v", schema.SchemaName, err)
			panic(err)
		}
		timber.Tracef("created schema [%s - %d]", s.SchemaName, s.SchemaID)
	}

	timber.Debugf("creating %d postgres table(s)", len(pgTables))
	for _, table := range pgTables {
		t, c, err := ctx.Tables().NewTable(table.Table, table.Columns)
		if err != nil {
			timber.Fatalf("failed to create table [%s]: %v", table.Table.TableName, err)
			panic(err)
		}
		timber.Tracef("created table [%s - %d] with %d column(s)", t.TableName, t.TableID, len(c))
	}
}
