package core

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
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

func TestSchemaContext_NewSchema_MultiServer(t *testing.T) {
	colony1, cleanup1 := newTestColony()
	defer cleanup1()

	colony2, cleanup2 := newTestColony(colony1.Addr().String())
	defer cleanup2()
	t.Run("create a new schema", func(t *testing.T) {
		name := "public"
		schema, err := colony1.Schema().(*schemaContext).NewSchema(name)
		assert.NoError(t, err)
		assert.NotEmpty(t, schema)
		assert.True(t, schema.SchemaID > 0)
		assert.Equal(t, name, schema.SchemaName)
		time.Sleep(1 * time.Second)
		exists, err := colony2.Schema().(*schemaContext).Exists(name)
		assert.NoError(t, err)
		assert.True(t, exists)
	})
}
