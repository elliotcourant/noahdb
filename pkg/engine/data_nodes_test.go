package engine_test

import (
	"github.com/elliotcourant/noahdb/pkg/engine"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDataNodeContextBase_NewDataNode(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
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
}

func TestDataNodeContextBase_GetDataNode(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		cluster, cleanup := NewTestCoreCluster(t, 1)
		defer cleanup()

		txn, err := cluster[0].Begin()
		assert.NoError(t, err)

		dataNode, err := txn.
			DataNodes().
			NewDataNode("postgres", 5432, "user", "password")
		assert.NoError(t, err)
		assert.NotZero(t, dataNode)

		retrievedDataNode, err := txn.DataNodes().GetDataNode(dataNode.DataNodeId)
		assert.NoError(t, err)
		assert.Equal(t, dataNode, retrievedDataNode)
	})

	// Make sure that if we pass a DataNodeId that does not exist that we receive an error.
	t.Run("not found", func(t *testing.T) {
		cluster, cleanup := NewTestCoreCluster(t, 1)
		defer cleanup()

		txn, err := cluster[0].Begin()
		assert.NoError(t, err)

		retrievedDataNode, err := txn.DataNodes().GetDataNode(1)
		assert.Equal(t, engine.ErrDataNodeNotFound, err)
		assert.Empty(t, retrievedDataNode)
	})
}

func TestDataNodeContextBase_GetDataNodes(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		numberOfDataNodes := 10

		cluster, cleanup := NewTestCoreCluster(t, 1)
		defer cleanup()

		txn, err := cluster[0].Begin()
		assert.NoError(t, err)

		for i := 0; i < numberOfDataNodes; i++ {
			dataNode, err := txn.
				DataNodes().
				NewDataNode("postgres", 5432+i, "user", "password")
			assert.NoError(t, err)
			assert.NotZero(t, dataNode)
		}

		dataNodes, err := txn.DataNodes().GetDataNodes()
		assert.NoError(t, err)
		assert.NotEmpty(t, dataNodes)
		assert.Len(t, dataNodes, numberOfDataNodes)
	})
}
