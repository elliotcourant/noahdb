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
