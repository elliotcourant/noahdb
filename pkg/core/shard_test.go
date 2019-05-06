package core

import (
	"github.com/readystock/golog"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestShardContext_NewShard(t *testing.T) {
	colony, cleanup := newTestColony()
	defer cleanup()
	newShard, err := colony.Shards().NewShard()
	assert.NoError(t, err)
	assert.True(t, newShard.ShardID > 0)
}

func TestShardContext_GetShards(t *testing.T) {
	colony, cleanup := newTestColony()
	defer cleanup()
	newShard, err := colony.Shards().NewShard()
	assert.NoError(t, err)
	assert.True(t, newShard.ShardID > 0)
	time.Sleep(1 * time.Second)
	golog.Verbosef("trying to query shards")
	shards, err := colony.Shards().GetShards()
	assert.NoError(t, err)
	assert.NotEmpty(t, shards)
}

func TestShardContext_GetWriteDataNodeShards(t *testing.T) {
	colony, cleanup := newTestColony()
	defer cleanup()
	shards, err := colony.Shards().GetWriteDataNodeShards(1)
	assert.NoError(t, err)
	assert.Empty(t, shards)
}

func TestShardContext_BalanceOrphanedShards(t *testing.T) {
	colony, cleanup := newTestColony()
	defer cleanup()
	t.Run("balance orphaned shards", func(t *testing.T) {
		newShard, err := colony.Shards().NewShard()
		assert.NoError(t, err)
		assert.True(t, newShard.ShardID > 0)
		err = colony.Shards().(*shardContext).BalanceOrphanShards()
		assert.NoError(t, err)
	})
}
