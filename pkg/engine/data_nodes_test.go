package engine_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDataNodeContextBase_NewDataNode(t *testing.T) {
	t.Run("single", func(t *testing.T) {
		cluster, cleanup := NewTestCoreCluster(t, 1)
		defer cleanup()

		txn, err := cluster[0].Begin()
		assert.NoError(t, err)

		dataNode, err := txn.
			DataNodes().
			NewDataNode("postgres", 5432, "user", "password")
		assert.NoError(t, err)
		assert.NotZero(t, dataNode)
	})

	t.Run("cluster", func(t *testing.T) {
		cluster, cleanup := NewTestCoreCluster(t, 3)
		defer cleanup()

		txn, err := cluster[0].Begin()
		assert.NoError(t, err)

		dataNode, err := txn.
			DataNodes().
			NewDataNode("postgres", 5432, "user", "password")
		assert.NoError(t, err)
		assert.NotZero(t, dataNode)
	})
}
