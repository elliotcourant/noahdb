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

		txn := cluster.Begin(t)

		dataNode, err := txn.
			DataNodes().
			NewDataNode("postgres", 5432, "user", "password")
		assert.NoError(t, err)
		assert.NotZero(t, dataNode)
	})

	t.Run("unique constraint", func(t *testing.T) {
		cluster, cleanup := NewTestCoreCluster(t, 1)
		defer cleanup()

		txn := cluster.Begin(t)

		_, err := txn.
			DataNodes().
			NewDataNode("postgres", 5432, "user", "password")
		assert.NoError(t, err)

		_, err = txn.
			DataNodes().
			NewDataNode("postgres", 5432, "user", "password")
		assert.Error(t, err)
	})

	t.Run("distributed unique constraint", func(t *testing.T) {
		cluster, cleanup := NewTestCoreCluster(t, 3)
		defer cleanup()

		txn1 := cluster.BeginOn(t, 0)
		txn2 := cluster.BeginOn(t, 1)

		_, err := txn1.
			DataNodes().
			NewDataNode("postgres", 5432, "user", "password")
		assert.NoError(t, err)

		_, err = txn2.
			DataNodes().
			NewDataNode("postgres", 5432, "user", "password")
		assert.NoError(t, err)

		err = txn1.Commit()
		assert.NoError(t, err)

		err = txn2.Commit()
		assert.Error(t, err)
	})
}

func TestDataNodeContextBase_GetDataNode(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		cluster, cleanup := NewTestCoreCluster(t, 1)
		defer cleanup()

		txn := cluster.Begin(t)

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

		txn := cluster.Begin(t)

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

		txn := cluster.Begin(t)

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

func TestDataNodeContextBase_GetDataNodesForShard(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		cluster, cleanup := NewTestCoreClusterEx(t, 1, true)
		defer cleanup()

		txn := cluster.Begin(t)

		dataNodes, err := txn.DataNodes().GetDataNodes()
		assert.NoError(t, err)

		_, err = txn.
			DataNodeShards().
			NewDataNodeShard(dataNodes[0].DataNodeId, 1, engine.DataNodeShardPosition_Leader)
		assert.NoError(t, err)

		_, err = txn.
			DataNodeShards().
			NewDataNodeShard(dataNodes[0].DataNodeId, 2, engine.DataNodeShardPosition_Leader)
		assert.NoError(t, err)

		dataNodes1, err := txn.DataNodes().GetDataNodesForShard(1, engine.DataNodeShardPosition_Leader)
		assert.NoError(t, err)
		assert.Equal(t, []engine.DataNode{dataNodes[0]}, dataNodes1)

		dataNodes2, err := txn.DataNodes().GetDataNodesForShard(2, engine.DataNodeShardPosition_Leader)
		assert.NoError(t, err)
		assert.Equal(t, []engine.DataNode{dataNodes[0]}, dataNodes2)
	})

	// Make sure that we can pass a position to the data node shard getter so we can get data nodes
	// for targeted queries.
	t.Run("filter by position", func(t *testing.T) {
		cluster, cleanup := NewTestCoreClusterEx(t, 2, true)
		defer cleanup()

		txn := cluster.Begin(t)

		dataNodes, err := txn.DataNodes().GetDataNodes()
		assert.NoError(t, err)

		_, err = txn.
			DataNodeShards().
			NewDataNodeShard(dataNodes[0].DataNodeId, 1, engine.DataNodeShardPosition_Leader)
		assert.NoError(t, err)

		_, err = txn.
			DataNodeShards().
			NewDataNodeShard(dataNodes[1].DataNodeId, 1, engine.DataNodeShardPosition_Follower)
		assert.NoError(t, err)

		dataNodes1, err := txn.DataNodes().GetDataNodesForShard(1, engine.DataNodeShardPosition_Leader)
		assert.NoError(t, err)
		assert.Equal(t, []engine.DataNode{dataNodes[0]}, dataNodes1)

		dataNodes2, err := txn.DataNodes().GetDataNodesForShard(1, engine.DataNodeShardPosition_Follower)
		assert.NoError(t, err)
		assert.Equal(t, []engine.DataNode{dataNodes[1]}, dataNodes2)
	})
}
