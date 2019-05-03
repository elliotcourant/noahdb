package core

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTableContext_GetTables(t *testing.T) {
	colony, cleanup := newTestColony()
	defer cleanup()
	t.Run("get imaginary tables", func(t *testing.T) {
		tables, err := colony.Tables().
			GetTables("users", "accounts")
		assert.NoError(t, err)
		assert.Empty(t, tables)
	})
}

func TestTableContext_GetTablesInSchema(t *testing.T) {
	colony, cleanup := newTestColony()
	defer cleanup()
	t.Run("get tables in imaginary schema", func(t *testing.T) {
		tables, err := colony.Tables().(*tableContext).
			GetTablesInSchema("imaginary", "test")
		assert.NoError(t, err)
		assert.Empty(t, tables)
	})
}
