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

	t.Run("distributed unique constraint", func(t *testing.T) {
		cluster, cleanup := NewTestCoreCluster(t, 3)
		defer cleanup()

		txn1, err := cluster[0].Begin()
		assert.NoError(t, err)

		txn2, err := cluster[1].Begin()
		assert.NoError(t, err)

		_, err = txn1.
			DataNodeShards().
			NewDataNodeShard(1, 1, engine.DataNodeShardPosition_Leader)
		assert.NoError(t, err)

		_, err = txn2.
			DataNodeShards().
			NewDataNodeShard(1, 1, engine.DataNodeShardPosition_Leader)
		assert.NoError(t, err)

		err = txn1.Commit()
		assert.NoError(t, err)

		err = txn2.Commit()
		assert.Error(t, err)
	})
}
