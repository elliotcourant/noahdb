package engine_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPoolContextBase_GetConnection(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		cluster, cleanup := NewTestCoreClusterEx(t, 1, true)
		defer cleanup()

		txn := cluster.Begin(t)

		shard, dataNodeShards, err := txn.Shards().NewShard()
		if !assert.NoError(t, err) {
			panic(err)
		}
		assert.NotZero(t, shard.ShardId)
		assert.NotEmpty(t, dataNodeShards)
		assert.Len(t, dataNodeShards, 1)

		dataNodeShard := dataNodeShards[0]

		conn, err := txn.Connect().GetConnection(dataNodeShard.DataNodeShardId)
		assert.NoError(t, err)
		assert.NotNil(t, conn)
	})
}
