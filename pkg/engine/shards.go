package engine

import (
	"errors"
	"github.com/ahmetb/go-linq/v3"
	"github.com/elliotcourant/mellivora"
)

var (
	// ErrShardNotFound is returned when a shard is requested by it's Id but a record with that Id
	// does not exist.
	ErrShardNotFound = errors.New("shard does not exist")

	// ErrNoDataNodesToPlaceShard is returned when a shard is created but there are no data nodes in
	// the cluster to actually place the shard on to.
	ErrNoDataNodesToPlaceShard = errors.New("no data nodes are available to place shard")
)

var (
	_ ShardContext = &shardBaseContext{}
)

type (
	// ShardState is an indicator for the current state of the shard (duh). It will indicate whether
	// or not a shard can be used at all within a cluster. In the future it will be used to indicate
	// if a shard is currently being rebalanced or if it's being restored.
	ShardState int

	// Shard keeps track of what shards are in the cluster and what their current state is.
	Shard struct {
		ShardId uint64 `m:"pk"`
		State   ShardState
	}

	// ShardContext provides an accessor interface for shard models.
	ShardContext interface {
		// GetShard will return a shard with the Id provided. If a shard does not exist with that Id
		// then an ErrShardNotFound error will be returned.
		GetShard(shardId uint64) (Shard, error)

		// GetShards will return all of the shards in the entire cluster. If no shard states are
		// provided, then all shard will be returned. Otherwise only shards that have one of the
		// shard states provided will be returned.
		GetShards(states ...ShardState) ([]Shard, error)

		// NewShard will create a new shard with the ShardState_Initializing state.
		NewShard() (Shard, []DataNodeShard, error)
	}

	shardBaseContext struct {
		t *transactionBase
	}
)

const (
	// ShardState_Unknown indicates that the state for the shard is missing or hasn't yet been
	// established. We want to ignore shards that have an unknown state.
	ShardState_Unknown ShardState = iota

	// ShardState_Initializing indicates that the shard has just been created and has not yet been
	// assigned to a set of data nodes. This shard cannot be used yet.
	ShardState_Initializing

	// ShardState_Ready indicates that the shard is healthy and is capable of being used for new
	// tenants or for other actions.
	ShardState_Ready

	// ShardState_ReadOnly indicates that a shard's replicas have fallen below the minimum required
	// to maintain a confidently writable state. This happens when half or more of a shards replicas
	// become unavailable.
	// TODO (elliotcourant) in the event of a split brain, the minority cluster will not be able to
	//  elect a leader to communicate to the other minority nodes that certain shards are read only.
	//  A method should be devised to determine if a shard is read only independent of this state.
	ShardState_ReadOnly

	// ShardState_Unavailable indicates that all of a shard's replicas are not reachable by the
	// cluster.
	ShardState_Unavailable
)

// Shards returns the accessors for shards.
func (t *transactionBase) Shards() ShardContext {
	return &shardBaseContext{
		t: t,
	}
}

// GetShard will return a shard with the Id provided. If a shard does not exist with that Id
// then an ErrShardNotFound error will be returned.
func (s *shardBaseContext) GetShard(shardId uint64) (Shard, error) {
	shard := Shard{}
	err := s.t.txn.Model(shard).Where(mellivora.Ex{
		"ShardId": shardId,
	}).Select(&shard)
	if shard.ShardId == 0 && err == nil {
		return shard, ErrShardNotFound
	}

	return shard, err
}

// GetShards will return all of the shards in the entire cluster. If no shard states are
// provided, then all shard will be returned. Otherwise only shards that have one of the
// shard states provided will be returned.
func (s *shardBaseContext) GetShards(states ...ShardState) ([]Shard, error) {
	shards := make([]Shard, 0)
	query := s.t.txn.Model(shards)

	// If there are states specified that we want to filter by make sure to restrict the results
	// to those states.
	if len(states) > 0 {
		query = query.Where(mellivora.Ex{
			"State": states,
		})
	}

	err := query.Select(&shards)

	return shards, err
}

// NewShard will create a new shard with the ShardState_Initializing state. The shard will
// remain in that state until the leader of the cluster picks the shard up and assigns it to
// it's data nodes.
func (s *shardBaseContext) NewShard() (Shard, []DataNodeShard, error) {
	shard := Shard{}

	shardDist, err := s.t.DataNodes().GetDataNodeShardDistribution()
	if err != nil {
		return shard, nil, err
	} else if len(shardDist) == 0 {
		return shard, nil, ErrNoDataNodesToPlaceShard
	}

	id, err := s.t.core.store.NextSequenceId("shards")
	if err != nil {
		return shard, nil, err
	}

	shard.ShardId = id

	// TODO (elliotcourant) move the replication factor variable to a setting.
	replicationFactor := 3

	if len(shardDist) < replicationFactor {
		// TODO (elliotcourant) add something here to handle gracefully when a shard MUST be
		//  under-replicated due to lack of nodes.
		replicationFactor = len(shardDist)
	}

	dataNodeIds := make([]uint64, 0)
	linq.From(shardDist).
		// Sort the data nodes by the number of shards they have in ascending order.
		OrderBy(func(i interface{}) interface{} {
			return i.(linq.KeyValue).Value
		}).
		// Take the top X data nodes with the least number of shards.
		Take(replicationFactor).
		// Pull the DataNodeId out.
		Select(func(i interface{}) interface{} {
			return i.(linq.KeyValue).Key
		}).
		ToSlice(&dataNodeIds)

	// We now know what data nodes we will be placing our shards on. Create the data node shard
	// records to keep track of which node will be which position.
	dataNodeShards := make([]DataNodeShard, 0)
	for i, dataNodeId := range dataNodeIds {
		// By default each one of the data node shards is a follower.
		position := DataNodeShardPosition_Follower

		// But the first data node shard that we create will be the leader.
		if i == 0 {
			position = DataNodeShardPosition_Leader
		}

		dataNodeShard, err := s.t.
			DataNodeShards().
			NewDataNodeShard(dataNodeId, shard.ShardId, position)
		if err != nil {
			return shard, nil, err
		}

		dataNodeShards = append(dataNodeShards, dataNodeShard)
	}

	shard.State = ShardState_Ready

	return shard, dataNodeShards, s.t.txn.Insert(shard)
}
