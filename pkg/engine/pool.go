package engine

import (
	"github.com/elliotcourant/noahdb/pkg/pgproto"
	"net"
)

var (
	_ PoolConnection = &dataNodeShardConnection{}
	_ PoolContext    = &poolContextBase{}
)

type (
	poolContextBase struct {
		t *transactionBase
	}

	dataNodeShardConnection struct {
		conn net.Conn
		pgproto.Frontend
	}

	// PoolConnection is an interface around a single connection to a single data node shard.
	PoolConnection interface {
		pgproto.Frontend
		Close() error
		Release()
		ID() uint64
		IsRoot() bool
	}

	PoolContext interface {
		// GetConnection will return a connection to the specific database that is hosting the
		// data node shard. Only that particular shard is accessible from this connection,
		GetConnection(dataNodeShardId uint64) (PoolConnection, error)

		// GetRootConnection will return a connection to the POSTGRES database for the specified
		// data node. This can be used to make configuration changes to the database or to manage
		// replication.
		GetRootConnection(dataNodeId uint64) (PoolConnection, error)
	}
)

// Pool will return the accessor interface for the coordinator's data node pool..
func (t *transactionBase) Pool() PoolContext {
	return &poolContextBase{
		t: t,
	}
}

// GetConnection will return a connection to the specific database that is hosting the
// data node shard. Only that particular shard is accessible from this connection,
func (p poolContextBase) GetConnection(dataNodeShardId uint64) (PoolConnection, error) {
	panic("implement me")
}

// GetRootConnection will return a connection to the POSTGRES database for the specified
// data node. This can be used to make configuration changes to the database or to manage
// replication.
func (p poolContextBase) GetRootConnection(dataNodeId uint64) (PoolConnection, error) {
	panic("implement me")
}

func (c *dataNodeShardConnection) ID() uint64 {
	return 0 // TODO (elliotcourant) return the id of the data node shard.
}

// Release will return the connection to the pool if the connection is still available.
func (c *dataNodeShardConnection) Release() {
	if c.Frontend == nil {
		return
	}

}

// Close will invalidate this connection and make it no longer usable, it will not be returned to
// the pool.
func (c *dataNodeShardConnection) Close() error {
	c.Frontend = nil
	return c.conn.Close()
}

// IsRoot will return true if the current pool connection is targeting a non shard database.
func (c *dataNodeShardConnection) IsRoot() bool {
	return false
}
