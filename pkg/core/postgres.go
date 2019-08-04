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
				TableType: TableType_Postgres,
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
				{
					ColumnName: "aggcombinefn",
					Type:       types.Type_regproc,
				},
				{
					ColumnName: "aggserialfn",
					Type:       types.Type_regproc,
				},
				{
					ColumnName: "aggdeserialfn",
					Type:       types.Type_regproc,
				},
				{
					ColumnName: "aggmtransfn",
					Type:       types.Type_regproc,
				},
				{
					ColumnName: "aggminvtransfn",
					Type:       types.Type_regproc,
				},
				{
					ColumnName: "aggmfinalfn",
					Type:       types.Type_regproc,
				},
				{
					ColumnName: "aggfinalextra",
					Type:       types.Type_bool,
				},
				{
					ColumnName: "aggmfinalextra",
					Type:       types.Type_bool,
				},
				{
					ColumnName: "aggfinalmodify",
					Type:       types.Type_char,
				},
				{
					ColumnName: "aggmfinalmodify",
					Type:       types.Type_char,
				},
				{
					ColumnName: "aggsortop",
					Type:       types.Type_oid,
				},
				{
					ColumnName: "aggtranstype",
					Type:       types.Type_oid,
				},
				{
					ColumnName: "aggtransspace",
					Type:       types.Type_int4,
				},
				{
					ColumnName: "aggmtranstype",
					Type:       types.Type_oid,
				},
				{
					ColumnName: "aggmtransspace",
					Type:       types.Type_int4,
				},
				{
					ColumnName: "agginitval",
					Type:       types.Type_text,
				},
				{
					ColumnName: "aggminitval",
					Type:       types.Type_text,
				},
			},
		},
		{
			Table: Table{
				TableName: "pg_am",
				TableType: TableType_Postgres,
			},
			Columns: []Column{
				{
					ColumnName: "oid",
					Type:       types.Type_oid,
				},
				{
					ColumnName: "amname",
					Type:       types.Type_text,
				},
				{
					ColumnName: "amhandler",
					Type:       types.Type_regproc,
				},
				{
					ColumnName: "amtype",
					Type:       types.Type_char,
				},
			},
		},
		{
			Table: Table{
				TableName: "pg_amop",
				TableType: TableType_Postgres,
			},
			Columns: []Column{
				{
					ColumnName: "oid",
					Type:       types.Type_oid,
				},
				{
					ColumnName: "amopfamily",
					Type:       types.Type_oid,
				},
				{
					ColumnName: "ampolefttype",
					Type:       types.Type_oid,
				},
				{
					ColumnName: "amoprighttype",
					Type:       types.Type_oid,
				},
				{
					ColumnName: "amopstrategy",
					Type:       types.Type_int2,
				},
				{
					ColumnName: "amoppurpose",
					Type:       types.Type_char,
				},
				{
					ColumnName: "amopopr",
					Type:       types.Type_oid,
				},
				{
					ColumnName: "amopmethod",
					Type:       types.Type_oid,
				},
				{
					ColumnName: "amopsortfamily",
					Type:       types.Type_oid,
				},
			},
		},
		{
			Table: Table{
				TableName: "pg_amproc",
				TableType: TableType_Postgres,
			},
			Columns: []Column{
				{
					ColumnName: "oid",
					Type:       types.Type_oid,
				},
				{
					ColumnName: "amprocfamily",
					Type:       types.Type_oid,
				},
				{
					ColumnName: "amproclefttype",
					Type:       types.Type_oid,
				},
				{
					ColumnName: "amprocrighttype",
					Type:       types.Type_oid,
				},
				{
					ColumnName: "amprocnum",
					Type:       types.Type_int2,
				},
				{
					ColumnName: "amproc",
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
