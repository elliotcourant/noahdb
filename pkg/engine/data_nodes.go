package engine

import (
	"github.com/elliotcourant/mellivora"
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

	DataNodeContext interface {
		NewDataNode(address string, port int, user, password string) (DataNode, error)
		GetDataNode(dataNodeId uint64) (DataNode, error)
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
// data node Id.
func (d *dataNodeContextBase) GetDataNode(dataNodeId uint64) (DataNode, error) {
	dataNode := DataNode{}
	err := d.t.txn.
		Model(dataNode).
		Where(mellivora.Ex{
			"DataNodeId": dataNodeId,
		}).
		Select(&dataNode)
	return dataNode, err
}
