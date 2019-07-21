package core_test

import (
	"github.com/elliotcourant/noahdb/testutils"
	"github.com/elliotcourant/timber"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestShardContext_NewShard(t *testing.T) {
	colony, cleanup := testutils.NewTestColony(t)
	defer cleanup()
	newShard, err := colony.Shards().NewShard()
	assert.NoError(t, err)
	assert.True(t, newShard.ShardID > 0)
}

func TestShardContext_GetShards(t *testing.T) {
	colony, cleanup := testutils.NewTestColony(t)
	defer cleanup()
	newShard, err := colony.Shards().NewShard()
	assert.NoError(t, err)
	assert.True(t, newShard.ShardID > 0)
	time.Sleep(1 * time.Second)
	timber.Verbosef("trying to query shards")
	shards, err := colony.Shards().GetShards()
	assert.NoError(t, err)
	assert.NotEmpty(t, shards)
}

func TestShardContext_GetWriteDataNodeShards(t *testing.T) {
	colony, cleanup := testutils.NewPgTestColony(t)
	defer cleanup()
	shards, err := colony.Shards().GetWriteDataNodeShards(1)
	assert.NoError(t, err)
	assert.NotEmpty(t, shards)
}

func TestShardContext_BalanceOrphanedShards(t *testing.T) {
	t.Run("balance orphaned shards", func(t *testing.T) {
		colony, cleanup := testutils.NewPgTestColony(t)
		defer cleanup()
		newShard, err := colony.Shards().NewShard()
		assert.NoError(t, err)
		assert.True(t, newShard.ShardID > 0)
		err = colony.Shards().BalanceOrphanShards()
		assert.NoError(t, err)
	})

	t.Run("balance multiple orphaned shards", func(t *testing.T) {
		colony, cleanup := testutils.NewPgTestColony(t)
		defer cleanup()

		existingNodes, err := colony.DataNodes().GetDataNodes()
		assert.NoError(t, err)
		assert.NotEmpty(t, existingNodes)

		numberOfNodes := 5 - len(existingNodes)
		numberOfShards := 32

		for i := 0; i < numberOfNodes; i++ {
			node, cleanup, err := testutils.NewDataNode(t)
			if !assert.NoError(t, err) {
				panic(err)
			}
			defer cleanup()
			// Noahdb does not check to see if a node is a duplicate. And shards on each "node"
			// are unique so a single node can be treated as multiple nodes in somne cases.
			_, err = colony.DataNodes().NewDataNode(node.Address, node.Port, node.User, node.Password)
			if !assert.NoError(t, err) {
				panic(err)
			}
		}

		for i := 0; i < numberOfShards; i++ {
			_, _ = colony.Shards().NewShard()
		}

		pressureBefore, _ := colony.Shards().GetDataNodesPressure(numberOfNodes)
		for _, pressure := range pressureBefore {
			assert.Empty(t, pressure.Shards)
		}

		err = colony.Shards().BalanceOrphanShards()
		assert.NoError(t, err)

		pressureAfter, _ := colony.Shards().GetDataNodesPressure(numberOfNodes)
		for _, pressure := range pressureAfter {
			assert.NotEmpty(t, pressure.Shards)
		}
	})
}
