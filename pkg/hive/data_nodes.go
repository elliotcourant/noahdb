package hive

import (
	"encoding/binary"
	"github.com/elliotcourant/noahdb/pkg/core"
	"github.com/elliotcourant/timber"
	"github.com/gogo/protobuf/proto"
	"time"
)

type DataNodeContext interface {
	AddDataNode(dataNode core.DataNode) (core.DataNode, error)
	GetDataNodes() ([]core.DataNode, error)
	GetDataNode(id uint64) (dataNode core.DataNode, ok bool, err error)
}

func (h *hiveTransaction) DataNodes() DataNodeContext {
	return &dataNodesContext{
		txn: h,
	}
}

type dataNodesContext struct {
	txn *hiveTransaction
}

func (d *dataNodesContext) AddDataNode(dataNode core.DataNode) (core.DataNode, error) {
	id, err := d.txn.txn.NextIncrementId(dataNodesPrefix.Bytes())
	if err != nil {
		return dataNode, err
	}
	dataNode.DataNodeID = id
	path := append(dataNodesPrefix.Bytes(), make([]byte, 8)...)
	binary.BigEndian.PutUint64(path[1:], id)

	val, err := proto.Marshal(&dataNode)
	if err != nil {
		return dataNode, err
	}

	return dataNode, d.txn.txn.Set(path, val)
}

func (d *dataNodesContext) GetDataNodes() ([]core.DataNode, error) {
	start := time.Now()
	defer timber.Tracef("get data nodes time: %s", time.Since(start))
	dataNodes := make([]core.DataNode, 0)
	itr := d.txn.txn.GetIterator(dataNodesPrefix.Bytes(), false, false)
	defer itr.Close()
	for itr.Seek(dataNodesPrefix.Bytes()); itr.ValidForPrefix(dataNodesPrefix.Bytes()); itr.Next() {
		dataNode := core.DataNode{}
		if err := itr.Item().Value(func(val []byte) error {
			return proto.Unmarshal(val, &dataNode)
		}); err != nil {
			return nil, err
		}
		dataNodes = append(dataNodes, dataNode)
	}
	return dataNodes, nil
}

func (d *dataNodesContext) GetDataNode(id uint64) (dataNode core.DataNode, ok bool, err error) {
	start := time.Now()
	defer timber.Tracef("get data node time: %s", time.Since(start))
	path := append(dataNodesPrefix.Bytes(), make([]byte, 8)...)
	binary.BigEndian.PutUint64(path[1:], id)
	item := make([]byte, 0)
	item, ok, err = d.txn.txn.Get(path)
	if err != nil || !ok {
		return dataNode, ok, err
	}
	err = proto.Unmarshal(item, &dataNode)
	return
}
