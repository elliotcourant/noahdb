package engine

import (
	"errors"
	"github.com/elliotcourant/mellivora"
)

var (
	// ErrDataNodeNotFound is returned when the user requests a single
	// specific data node by Id but one does not exist.
	ErrDataNodeNotFound = errors.New("data node does not exist")
)

var (
	_ DataNodeContext = &dataNodeContextBase{}
)

type (
	// DataNode represents a reachable PostgreSQL database that is managed by the NoahDB cluster.
	DataNode struct {
		DataNodeId uint64 `m:"pk"`
		Address    string `m:"uq:uq_address_port"`
		Port       int    `m:"uq:uq_address_port"`
		Username   string
		Password   string
	}

	// DataNodeContext provides an accessor interface for data node models.
	DataNodeContext interface {
		// NewDataNode will create a data node object to allow shards to be placed on this node.
		// An error can be returned if there is a unique constraint violation.
		NewDataNode(address string, port int, user, password string) (DataNode, error)

		// GetDataNode will return a single data node struct that has the matching
		// data node Id. If no data node is found with the Id specified then
		// ErrDataNodeNotFound will be returned.
		GetDataNode(dataNodeId uint64) (DataNode, error)

		// GetDataNodes will return all of the data nodes in the entire cluster.
		GetDataNodes() ([]DataNode, error)

		// GetDataNodesForShard will return all of the data nodes that host a given shardId.
		GetDataNodesForShard(shardId uint64, positions ...DataNodeShardPosition) ([]DataNode, error)

		// GetDataNodeShardDistribution will return a map of data node Id's with counts of how many
		// shards are located on each data node.
		GetDataNodeShardDistribution() (map[uint64]int, error)
	}

	dataNodeContextBase struct {
		t *transactionBase
	}
)

// DataNodes returns accessors for data nodes.
func (t *transactionBase) DataNodes() DataNodeContext {
	return &dataNodeContextBase{
		t: t,
	}
}

// NewDataNode will create a data node object to allow shards to be placed on this node.
// An error can be returned if there is a unique constraint violation.
func (d *dataNodeContextBase) NewDataNode(address string, port int, user, password string) (DataNode, error) {
	dataNode := DataNode{}
	id, err := d.t.core.store.NextSequenceId("dataNodes")
	if err != nil {
		return dataNode, err
	}

	dataNode.DataNodeId = id
	dataNode.Address = address
	dataNode.Port = port
	dataNode.Username = user
	dataNode.Password = password

	return dataNode, d.t.txn.Insert(dataNode)
}

// GetDataNode will return a single data node struct that has the matching
// data node Id. If no data node is found with the Id specified then
// ErrDataNodeNotFound will be returned.
func (d *dataNodeContextBase) GetDataNode(dataNodeId uint64) (DataNode, error) {
	dataNode := DataNode{}
	err := d.t.txn.
		Model(dataNode).
		Where(mellivora.Ex{
			"DataNodeId": dataNodeId,
		}).
		Select(&dataNode)
	if dataNode.DataNodeId == 0 && err == nil {
		return dataNode, ErrDataNodeNotFound
	}

	return dataNode, err
}

// GetDataNodes will return all of the data nodes in the entire cluster.
func (d *dataNodeContextBase) GetDataNodes() ([]DataNode, error) {
	dataNodes := make([]DataNode, 0)
	err := d.t.txn.
		Model(dataNodes).
		Select(&dataNodes)
	return dataNodes, err
}

// GetDataNodesForShard will return all of the data nodes that host a given shardId.
func (d *dataNodeContextBase) GetDataNodesForShard(shardId uint64, positions ...DataNodeShardPosition) ([]DataNode, error) {
	dataNodeShards := make([]DataNodeShard, 0)

	dataNodeShardQuery := d.t.txn.
		Model(dataNodeShards).
		Where(mellivora.Ex{
			"ShardId": shardId,
		})

	// If we do want to filter by a position then add the filter, otherwise don't
	if len(positions) > 0 {
		dataNodeShardQuery = dataNodeShardQuery.AndWhere(mellivora.Ex{
			"Position": positions,
		})
	}

	if err := dataNodeShardQuery.
		Select(&dataNodeShards); err != nil {
		return nil, err
	}

	dataNodeIds := make([]uint64, 0)
	for _, dataNodeShard := range dataNodeShards {
		dataNodeIds = append(dataNodeIds, dataNodeShard.DataNodeId)
	}

	dataNodes := make([]DataNode, 0)
	err := d.t.txn.
		Model(dataNodes).
		Where(mellivora.Ex{
			"DataNodeId": dataNodeIds,
		}).
		Select(&dataNodes)
	return dataNodes, err
}

// GetDataNodeShardDistribution will return a map of data node Id's with counts of how many
// shards are located on each data node.
func (d *dataNodeContextBase) GetDataNodeShardDistribution() (map[uint64]int, error) {
	dataNodes, err := d.GetDataNodes()
	if err != nil {
		return nil, err
	}

	distribution := map[uint64]int{}

	for _, dataNode := range dataNodes {
		distribution[dataNode.DataNodeId] = 0
	}

	dataNodeShards, err := d.t.DataNodeShards().GetDataNodeShards()
	if err != nil {
		return nil, err
	}

	for _, dataNodeShard := range dataNodeShards {
		distribution[dataNodeShard.DataNodeId] = distribution[dataNodeShard.DataNodeId] + 1
	}

	return distribution, nil
}
