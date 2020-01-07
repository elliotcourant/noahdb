package engine_test

import (
	"github.com/elliotcourant/noahdb/pkg/engine"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDataNodeShardContextBase_NewDataNodeShard(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		cluster, cleanup := NewTestCoreClusterEx(t, 1, true)
		defer cleanup()

		txn := cluster.Begin(t)

		dataNode, err := txn.
			DataNodeShards().
			NewDataNodeShard(1, 1, engine.DataNodeShardPosition_Leader)
		assert.NoError(t, err)
		assert.NotZero(t, dataNode)
	})

	t.Run("data node does not exist", func(t *testing.T) {
		cluster, cleanup := NewTestCoreCluster(t, 1)
		defer cleanup()

		txn := cluster.Begin(t)

		_, err := txn.
			DataNodeShards().
			NewDataNodeShard(2, 1, engine.DataNodeShardPosition_Leader)
		assert.Equal(t, engine.ErrDataNodeNotFound, err)
	})

	t.Run("unique constraint", func(t *testing.T) {
		cluster, cleanup := NewTestCoreClusterEx(t, 1, true)
		defer cleanup()

		txn := cluster.Begin(t)

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
		cluster, cleanup := NewTestCoreClusterEx(t, 3, true)
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

func TestDataNodeShardContextBase_GetDataNodeShard(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		cluster, cleanup := NewTestCoreClusterEx(t, 1, true)
		defer cleanup()

		txn := cluster.Begin(t)

		dataNodeShard, err := txn.
			DataNodeShards().
			NewDataNodeShard(1, 1, engine.DataNodeShardPosition_Leader)
		assert.NoError(t, err)
		assert.NotZero(t, dataNodeShard)

		result, ok, err := txn.DataNodeShards().GetDataNodeShard(dataNodeShard.DataNodeShardId)
		assert.Equal(t, dataNodeShard, result)
		assert.True(t, ok)
		assert.NoError(t, err)
	})
}
