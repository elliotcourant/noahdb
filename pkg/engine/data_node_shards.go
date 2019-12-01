package engine

var (
	_ DataNodeShardContext = &dataNodeContextBase{}
)

type (
	// DataNodeShardPosition indicates the role a particular data node/shard pair play within the
	// cluster.
	DataNodeShardPosition int

	// DataNodeShard keeps track of what shards are stored on what data nodes, and what behavior
	// a particular data node/shard pair should perform.
	DataNodeShard struct {
		DataNodeShardId uint64 `m:"pk"`
		DataNodeId      uint64 `m:"uq:uq_data_node_id_shard_id"`
		ShardId         uint64 `m:"uq:uq_data_node_id_shard_id"`
		Position        DataNodeShardPosition
	}

	// DataNodeShardContext provides an accessor interface for data node shard models.
	DataNodeShardContext interface {
		// GetDataNodeShards will return all of the data node/shard pairs in the entire cluster.
		GetDataNodeShards() ([]DataNodeShard, error)

		// NewDataNodeShard will create a new data node/shard pair to keep track of replication flow and
		// what shards are located on which data nodes.
		NewDataNodeShard(dataNodeId, shardId uint64, position DataNodeShardPosition) (DataNodeShard, error)
	}

	dataNodeShardContextBase struct {
		t *transactionBase
	}
)

const (
	// DataNodeShardPosition_Unknown indicates that the value for the position
	// is it's default. We want to have this to make sure we don't accidentally
	// assume a data node shard's position incorrectly. Thus it's value is 0.
	DataNodeShardPosition_Unknown DataNodeShardPosition = iota

	// DataNodeShardPosition_Leader indicates that the current data node/shard
	// pair is the leader for that particular shard. Other shards will receive
	// a logical replication feed from this data node/shard.
	DataNodeShardPosition_Leader

	// DataNodeShardPosition_Follower indicates that the current data node/shard
	// pair is among the followers for a particular shard. It is read only. But
	// since it is a follower, it is automatically a candidate for promotion to
	// leader IF the current leader fails.
	DataNodeShardPosition_Follower
)

// DataNodeShards returns the accessors for data node shards.
func (t *transactionBase) DataNodeShards() DataNodeShardContext {
	return &dataNodeContextBase{
		t: t,
	}
}

// GetDataNodeShards will return all of the data node/shard pairs in the entire cluster.
func (d *dataNodeContextBase) GetDataNodeShards() ([]DataNodeShard, error) {
	dataNodeShards := make([]DataNodeShard, 0)
	err := d.t.txn.Model(dataNodeShards).Select(&dataNodeShards)

	return dataNodeShards, err
}

// NewDataNodeShard will create a new data node/shard pair to keep track of replication flow and
// what shards are located on which data nodes.
func (d *dataNodeContextBase) NewDataNodeShard(dataNodeId, shardId uint64, position DataNodeShardPosition) (DataNodeShard, error) {
	dataNodeShard := DataNodeShard{}
	id, err := d.t.core.store.NextSequenceId("dataNodeShards")
	if err != nil {
		return dataNodeShard, err
	}

	dataNodeShard.DataNodeShardId = id
	dataNodeShard.DataNodeId = dataNodeId
	dataNodeShard.ShardId = shardId
	dataNodeShard.Position = position

	return dataNodeShard, d.t.txn.Insert(dataNodeShard)
}
