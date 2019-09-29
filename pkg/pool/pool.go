package pool

import (
	"github.com/elliotcourant/noahdb/pkg/core"
	"github.com/elliotcourant/noahdb/pkg/pgproto"
	"github.com/elliotcourant/timber"
	"sync"
	"time"
)

// Connection represents a tunnel to and from a data node.
// It is a wrapper around the core.PoolConnection interface
// but when this connection is released it is returned to the
// executor's connection pool instead.
type Connection interface {
	Send(pgproto.FrontendMessage) error
	Receive() (pgproto.BackendMessage, error)
	ID() uint64
}

func NewConnection(conn core.PoolConnection, pool *BasePool) Connection {
	return &BaseConnection{
		conn: conn,
		pool: pool,
	}
}

type BaseConnection struct {
	conn core.PoolConnection
	pool *BasePool
}

func (b *BaseConnection) Send(pgproto.FrontendMessage) error {
	panic("implement me")
}

func (b *BaseConnection) Receive() (pgproto.BackendMessage, error) {
	panic("implement me")
}

func (b *BaseConnection) ID() uint64 {
	panic("implement me")
}

// Pool is connection manager for the executor.
type Pool interface {
	GetConnection(dataNodeShardId uint64) (Connection, error)
}

func NewPool(colony core.Colony) Pool {
	return &BasePool{
		colony:   colony,
		poolSync: sync.Mutex{},
		pool:     map[uint64]Connection{},
	}
}

// BasePool implements the Pool interface.
type BasePool struct {
	log      timber.Logger
	colony   core.Colony
	poolSync sync.Mutex
	pool     map[uint64]Connection
}

// GetConnection returns a connection that is allocated
// to this executor's connection pool.
func (p *BasePool) GetConnection(dataNodeShardId uint64) (Connection, error) {
	start := time.Now()
	defer p.log.Verbosef("[%s] acquisition of connection to data node shard [%d]",
		time.Since(start),
		dataNodeShardId)
	p.poolSync.Lock()
	defer p.poolSync.Unlock()
	if pool, ok := p.pool[dataNodeShardId]; ok {
		return pool, nil
	}
	pc, err := p.colony.Pool().GetConnectionForDataNodeShard(dataNodeShardId)
	if err != nil {
		return nil, err
	}
	p.pool[dataNodeShardId] = NewConnection(pc, p)
	panic("implement me")
}
