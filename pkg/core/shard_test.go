package core_test

import (
	"github.com/elliotcourant/noahdb/testutils"
	"github.com/readystock/golog"
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
	golog.Verbosef("trying to query shards")
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

		templateNode := existingNodes[0]

		numberOfNodes := 10 - len(existingNodes)
		numberOfShards := 32

		for i := 0; i < numberOfNodes; i++ {
			// Noahdb does not check to see if a node is a duplicate. And shards on each "node"
			// are unique so a single node can be treated as multiple nodes in somne cases.
			_, _ = colony.DataNodes().NewDataNode(templateNode.Address, int(templateNode.Port), templateNode.User, templateNode.Password)
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
