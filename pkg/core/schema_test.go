package core

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSchemaContext_Exists(t *testing.T) {
	colony, cleanup := newTestColony()
	defer cleanup()
	t.Run("doesn't exist", func(t *testing.T) {
		ok, err := colony.Schema().(*schemaContext).Exists("imaginary")
		assert.NoError(t, err)
		assert.False(t, ok)
	})
}

func TestSchemaContext_NewSchema(t *testing.T) {
	colony, cleanup := newTestColony()
	defer cleanup()
	t.Run("create a new schema", func(t *testing.T) {
		name := "public"
		schema, err := colony.Schema().(*schemaContext).NewSchema(name)
		assert.NoError(t, err)
		assert.NotEmpty(t, schema)
		assert.True(t, schema.SchemaID > 0)
		assert.Equal(t, name, schema.SchemaName)
	})
}
