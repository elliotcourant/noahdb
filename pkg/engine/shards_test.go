package engine_test

import (
	"github.com/elliotcourant/noahdb/pkg/engine"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestShardBaseContext_NewShard(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		cluster, cleanup := NewTestCoreClusterEx(t, 1, true)
		defer cleanup()

		txn := cluster.Begin(t)

		dataNodes, err := txn.DataNodes().GetDataNodes()
		assert.NoError(t, err)

		shard, dataNodeShards, err := txn.Shards().NewShard()
		assert.NoError(t, err)
		assert.NotZero(t, shard.ShardId)
		assert.Equal(t, engine.ShardState_Ready, shard.State)
		assert.NotEmpty(t, dataNodeShards)
		assert.Equal(t, dataNodes[0].DataNodeId, dataNodeShards[0].DataNodeId)
	})

	t.Run("no data nodes", func(t *testing.T) {
		cluster, cleanup := NewTestCoreCluster(t, 1)
		defer cleanup()

		txn := cluster.Begin(t)

		shard, dataNodeShards, err := txn.Shards().NewShard()
		assert.Equal(t, engine.ErrNoDataNodesToPlaceShard, err)
		assert.Zero(t, shard.ShardId)
		assert.Nil(t, dataNodeShards)
	})

	t.Run("balance", func(t *testing.T) {
		numberOfDataNodes, numberOfShards := 9, 32

		cluster, cleanup := NewTestCoreCluster(t, 1)
		defer cleanup()

		cluster.SeedDataNodes(t, numberOfDataNodes)

		txn := cluster.Begin(t)

		for i := 0; i < numberOfShards; i++ {
			shard, dataNodeShards, err := txn.Shards().NewShard()
			assert.NoError(t, err)
			assert.NotEmpty(t, dataNodeShards)
			assert.Len(t, dataNodeShards, 3)
			assert.NotZero(t, shard.ShardId)
		}

		dist, err := txn.DataNodes().GetDataNodeShardDistribution()
		assert.NoError(t, err)
		for dataNodeId, shards := range dist {
			assert.NotZero(t, shards, "data node [%d] has 0 shards", dataNodeId)
		}
	})
}

func TestShardBaseContext_GetShard(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		cluster, cleanup := NewTestCoreClusterEx(t, 1, true)
		defer cleanup()

		txn := cluster.Begin(t)

		shard, _, err := txn.Shards().NewShard()
		assert.NoError(t, err)
		assert.NotZero(t, shard.ShardId)
		assert.Equal(t, engine.ShardState_Ready, shard.State)

		shardRead, err := txn.Shards().GetShard(shard.ShardId)
		assert.NoError(t, err)
		assert.Equal(t, shard, shardRead)
	})

	t.Run("not found", func(t *testing.T) {
		cluster, cleanup := NewTestCoreCluster(t, 1)
		defer cleanup()

		txn := cluster.Begin(t)

		_, err := txn.Shards().GetShard(1)
		assert.Equal(t, engine.ErrShardNotFound, err)
	})
}

func TestShardBaseContext_GetShards(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		numberOfShards := 10

		cluster, cleanup := NewTestCoreClusterEx(t, 1, true)
		defer cleanup()

		txn := cluster.Begin(t)

		for i := 0; i < numberOfShards; i++ {
			shard, _, err := txn.Shards().NewShard()
			assert.NoError(t, err)
			assert.NotZero(t, shard.ShardId)
			assert.Equal(t, engine.ShardState_Ready, shard.State)
		}

		shards, err := txn.Shards().GetShards()
		assert.NoError(t, err)
		assert.Len(t, shards, numberOfShards)

		shards, err = txn.Shards().GetShards(engine.ShardState_Ready)
		assert.NoError(t, err)
		assert.Len(t, shards, numberOfShards)

		shards, err = txn.Shards().GetShards(engine.ShardState_Unavailable)
		assert.NoError(t, err)
		assert.Len(t, shards, 0)
	})
}
