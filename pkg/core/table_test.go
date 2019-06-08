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
