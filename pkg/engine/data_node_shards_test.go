package engine_test

import (
	"github.com/elliotcourant/noahdb/pkg/engine"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDataNodeContextBase_NewDataNodeShard(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		cluster, cleanup := NewTestCoreCluster(t, 1)
		defer cleanup()

		txn, err := cluster[0].Begin()
		assert.NoError(t, err)

		dataNode, err := txn.
			DataNodeShards().
			NewDataNodeShard(1, 1, engine.DataNodeShardPosition_Leader)
		assert.NoError(t, err)
		assert.NotZero(t, dataNode)
	})

	t.Run("unique constraint", func(t *testing.T) {
		cluster, cleanup := NewTestCoreCluster(t, 1)
		defer cleanup()

		txn, err := cluster[0].Begin()
		assert.NoError(t, err)

		dataNode, err := txn.
			DataNodeShards().
			NewDataNodeShard(1, 1, engine.DataNodeShardPosition_Leader)
		assert.NoError(t, err)
		assert.NotZero(t, dataNode)

		_, err = txn.
			DataNodeShards().
			NewDataNodeShard(1, 1, engine.DataNodeShardPosition_Leader)
		assert.Error(t, err)
	})
}
