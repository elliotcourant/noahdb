package pool

import (
	"github.com/elliotcourant/noahdb/pkg/core"
	"github.com/elliotcourant/timber"
	"sync"
	"time"
)

// Pool is connection manager for the executor.
type Pool interface {
	Begin()
	GetConnection(dataNodeShardId uint64, txn bool) (Connection, error)
}

func NewPool(colony core.Colony, logger timber.Logger) Pool {
	return &BasePool{
		colony:   colony,
		poolSync: sync.Mutex{},
		pool:     map[uint64]Connection{},
		log:      logger,
	}
}

// BasePool implements the Pool interface.
type BasePool struct {
	log           timber.Logger
	colony        core.Colony
	poolSync      sync.Mutex
	pool          map[uint64]Connection
	inTransaction bool
}

func (p *BasePool) Begin() {
	p.inTransaction = true
}

// GetConnection returns a connection that is allocated
// to this executor's connection pool.
func (p *BasePool) GetConnection(dataNodeShardId uint64, txn bool) (Connection, error) {
	start := time.Now()
	defer p.log.Verbosef("[%s] acquisition of connection to data node shard [%d]",
		time.Since(start),
		dataNodeShardId)
	p.poolSync.Lock()
	defer p.poolSync.Unlock()
	if pool, ok := p.pool[dataNodeShardId]; ok {
		// If txn is true then we need to make sure that this connection is
		// transactional.
		if txn {
			if err := pool.Begin(); err != nil {
				return nil, err
			}
		}
		return pool, nil
	}
	var conn Connection
	if pc, err := p.colony.Pool().GetConnectionForDataNodeShard(dataNodeShardId); err != nil {
		return nil, err
	} else {
		conn = NewConnection(pc, p)
	}

	p.pool[dataNodeShardId] = conn
	// If the current pool is in a transaction, or if we are executing this statement
	// transactionally then we need to send begin.
	if p.inTransaction || txn {
		if err := conn.Begin(); err != nil {
			return nil, err
		}
	}
	return conn, nil
}
