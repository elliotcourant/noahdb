package core_test

import (
	"github.com/elliotcourant/noahdb/testutils"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSchemaContext_Exists(t *testing.T) {
	colony, cleanup := testutils.NewTestColony()
	defer cleanup()
	t.Run("doesn't exist", func(t *testing.T) {
		ok, err := colony.Schema().Exists("imaginary")
		assert.NoError(t, err)
		assert.False(t, ok)
	})
}

func TestSchemaContext_NewSchema(t *testing.T) {
	colony, cleanup := testutils.NewTestColony()
	defer cleanup()
	t.Run("create a new schema", func(t *testing.T) {
		name := "public"
		schema, err := colony.Schema().NewSchema(name)
		assert.NoError(t, err)
		assert.NotEmpty(t, schema)
		assert.True(t, schema.SchemaID > 0)
		assert.Equal(t, name, schema.SchemaName)
	})
}

func TestSchemaContext_NewSchema_MultiServer(t *testing.T) {
	t.Run("create a new schema", func(t *testing.T) {
		colony1, cleanup1 := testutils.NewTestColony()
		defer cleanup1()
		name := "public"
		schema, err := colony1.Schema().NewSchema(name)
		assert.NoError(t, err)
		assert.NotEmpty(t, schema)
		assert.True(t, schema.SchemaID > 0)
		assert.Equal(t, name, schema.SchemaName)
		colony2, cleanup2 := testutils.NewTestColony(colony1.Addr().String())
		defer cleanup2()
		time.Sleep(1 * time.Second)
		exists, err := colony2.Schema().Exists(name)
		assert.NoError(t, err)
		assert.True(t, exists)
	})
}
