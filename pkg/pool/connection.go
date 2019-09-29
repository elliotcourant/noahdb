package pool

import (
	"fmt"
	"github.com/elliotcourant/noahdb/pkg/core"
	"github.com/elliotcourant/noahdb/pkg/pgproto"
	"github.com/elliotcourant/timber"
)

// Connection represents a tunnel to and from a data node.
// It is a wrapper around the core.PoolConnection interface
// but when this connection is released it is returned to the
// executor's connection pool instead.
type Connection interface {
	Begin() error
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
	log  timber.Logger
	conn core.PoolConnection
	pool *BasePool
}

func (b *BaseConnection) Begin() error {
	return b.Exec("BEGIN")
}

func (b *BaseConnection) Commit() error {
	return b.Exec("COMMIT")
}

func (b *BaseConnection) Rollback() error {
	return b.Exec("ROLLBACK")
}

func (b *BaseConnection) Send(msg pgproto.FrontendMessage) error {
	return b.conn.Send(msg)
}

func (b *BaseConnection) Receive() (pgproto.BackendMessage, error) {
	return b.conn.Receive()
}

func (b *BaseConnection) ID() uint64 {
	return b.conn.ID()
}

func (b *BaseConnection) Exec(query string) error {
	if err := b.Send(&pgproto.Query{
		String: query,
	}); err != nil {
		return err
	}
	return b.waitForReady()
}

func (b *BaseConnection) waitForReady() error {
	for {
		msg, err := b.Receive()
		if err != nil {
			return err
		}
		switch m := msg.(type) {
		case *pgproto.ErrorResponse:
			return fmt.Errorf("received error from pool conn with begin: %v", m.Message)
		case *pgproto.ReadyForQuery:
			return nil
		default:
			b.log.Warningf("received unexpected message type [%T]", m)
			panic("received unexpected message type")
		}
	}
}
