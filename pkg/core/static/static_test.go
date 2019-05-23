package static

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetEmbeddedFile(t *testing.T) {
	t.Run("existing file", func(t *testing.T) {
		data, err := GetEmbeddedFile("/00_internal_sql.sql")
		assert.NoError(t, err)
		assert.NotEmpty(t, data)
	})
	t.Run("non-existant file", func(t *testing.T) {
		data, err := GetEmbeddedFile("/kjdsahjkdsa.asdjkas")
		assert.Error(t, err)
		assert.Empty(t, data)
	})
}
