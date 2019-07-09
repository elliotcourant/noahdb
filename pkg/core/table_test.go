package core_test

import (
	"github.com/elliotcourant/noahdb/pkg/core"
	"github.com/elliotcourant/noahdb/testutils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTableContext_NewTable(t *testing.T) {
	colony, cleanup := testutils.NewTestColony(t)
	defer cleanup()
	t.Run("get imaginary tables", func(t *testing.T) {
		table := core.Table{
			TableName: "accounts",
			TableType: core.TableType_Tenant,
		}
		columns := []core.Column{
			{
				ColumnName: "id",
				Type:       core.Type_int8,
				PrimaryKey: true,
			},
			{
				ColumnName: "name",
				Type:       core.Type_text,
				PrimaryKey: false,
				Nullable:   true,
			},
		}
		newTable, newColumns, err := colony.Tables().NewTable(table, columns)
		assert.NoError(t, err)
		assert.True(t, newTable.TableID > 0)
		assert.NotEmpty(t, newColumns)
	})
}

func TestTableContext_GetTables(t *testing.T) {
	colony, cleanup := testutils.NewTestColony(t)
	defer cleanup()
	t.Run("get imaginary tables", func(t *testing.T) {
		tables, err := colony.Tables().
			GetTables("users", "accounts")
		assert.NoError(t, err)
		assert.Empty(t, tables)
	})
}

func TestTableContext_GetTablesInSchema(t *testing.T) {
	colony, cleanup := testutils.NewTestColony(t)
	defer cleanup()
	t.Run("get tables in imaginary schema", func(t *testing.T) {
		tables, err := colony.Tables().GetTablesInSchema("imaginary", "test")
		assert.NoError(t, err)
		assert.Empty(t, tables)
	})
}

func TestTableContext_GetColumnFromTables(t *testing.T) {
	colony, cleanup := testutils.NewTestColony(t)
	defer cleanup()
	t.Run("get column from tables", func(t *testing.T) {
		table1 := core.Table{
			TableName: "accounts",
			TableType: core.TableType_Tenant,
		}
		columns1 := []core.Column{
			{
				ColumnName: "id",
				Type:       core.Type_int8,
				PrimaryKey: true,
			},
			{
				ColumnName: "name",
				Type:       core.Type_text,
				PrimaryKey: false,
				Nullable:   true,
			},
		}
		newTable1, newColumns1, err := colony.Tables().NewTable(table1, columns1)
		assert.NoError(t, err)
		assert.True(t, newTable1.TableID > 0)
		assert.NotEmpty(t, newColumns1)

		table2 := core.Table{
			TableName: "users",
			TableType: core.TableType_Tenant,
		}
		columns2 := []core.Column{
			{
				ColumnName: "not_id",
				Type:       core.Type_int8,
				PrimaryKey: true,
			},
			{
				ColumnName: "different_name",
				Type:       core.Type_text,
				PrimaryKey: false,
				Nullable:   true,
			},
		}
		newTable2, newColumns2, err := colony.Tables().NewTable(table2, columns2)
		assert.NoError(t, err)
		assert.True(t, newTable2.TableID > 0)
		assert.NotEmpty(t, newColumns2)

		column, ok, err := colony.Tables().GetColumnFromTables("id", []string{
			"users",
			"accounts",
		})
		assert.NoError(t, err)
		assert.True(t, ok)
		assert.NotEmpty(t, column)
	})

	t.Run("get ambiguous column from tables", func(t *testing.T) {
		table1 := core.Table{
			TableName: "accounts_2",
			TableType: core.TableType_Tenant,
		}
		columns1 := []core.Column{
			{
				ColumnName: "id",
				Type:       core.Type_int8,
				PrimaryKey: true,
			},
			{
				ColumnName: "name",
				Type:       core.Type_text,
				PrimaryKey: false,
				Nullable:   true,
			},
		}
		newTable1, newColumns1, err := colony.Tables().NewTable(table1, columns1)
		assert.NoError(t, err)
		assert.True(t, newTable1.TableID > 0)
		assert.NotEmpty(t, newColumns1)

		table2 := core.Table{
			TableName: "users_2",
			TableType: core.TableType_Tenant,
		}
		columns2 := []core.Column{
			{
				ColumnName: "id",
				Type:       core.Type_int8,
				PrimaryKey: true,
			},
			{
				ColumnName: "different_name",
				Type:       core.Type_text,
				PrimaryKey: false,
				Nullable:   true,
			},
		}
		newTable2, newColumns2, err := colony.Tables().NewTable(table2, columns2)
		assert.NoError(t, err)
		assert.True(t, newTable2.TableID > 0)
		assert.NotEmpty(t, newColumns2)

		column, ok, err := colony.Tables().GetColumnFromTables("id", []string{
			"users_2",
			"accounts_2",
		})
		assert.Error(t, err)
		assert.False(t, ok)
		assert.Empty(t, column)
	})
}
