package engine

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/elliotcourant/noahdb/pkg/pgproto"
	"github.com/elliotcourant/timber"
	"net"
	"sync"
)

var (
	ErrDataNodeShardNotFound = errors.New("data node shard does not exist")
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
		dataNodeShardId uint64
		dataNodeId      uint64
		conn            net.Conn
		pgproto.Frontend
	}

	dataNodeShardPool struct {
		pool sync.Pool
	}

	// PoolConnection is an interface around a single connection to a single data node shard.
	PoolConnection interface {
		pgproto.Frontend
		Close() error
		Release()
		DataNodeID() uint64
		ShardID() uint64
		DataNodeShardID() uint64
	}

	PoolContext interface {
		// GetConnection will return a connection to the specific database that is hosting the
		// data node shard. Only that particular shard is accessible from this connection,
		GetConnection(dataNodeShardId uint64) (PoolConnection, error)
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
func (p *poolContextBase) GetConnection(dataNodeShardId uint64) (PoolConnection, error) {
	dataNodeShard, ok, err := p.t.DataNodeShards().GetDataNodeShard(dataNodeShardId)
	if err != nil {
		return nil, err
	} else if !ok {
		return nil, ErrDataNodeShardNotFound
	}

	dataNode, err := p.t.DataNodes().GetDataNode(dataNodeShard.DataNodeId)
	if err != nil {
		return nil, err
	}

	address := fmt.Sprintf("%s:%d", dataNode.Address, dataNode.Port)

	// TODO (elliotcourant) If we fail to establish a connection to the database after this point
	//  then that means the error was not due to connection, but was due to poor authentication or
	//  another error. If that is the case then we want to close this network connection.
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}

	frontend, err := pgproto.NewFrontend(conn, conn)
	if err != nil {
		return nil, err
	}

	connection := &dataNodeShardConnection{
		dataNodeShardId: dataNodeShardId,
		dataNodeId:      dataNode.DataNodeId,
		conn:            conn,
		Frontend:        frontend,
	}

	startupMessage := pgproto.StartupMessage{
		ProtocolVersion: pgproto.ProtocolVersionNumber,
		Parameters: map[string]string{
			"user":             dataNode.Username,
			"application_name": fmt.Sprintf("noahdb_%s", p.t.core.store.NodeID()),
			"database":         getPgDatabaseName(dataNodeShardId),
			"client_encoding":  "UTF8",
		},
	}

	if err := connection.Send(&startupMessage); err != nil {
		return nil, err
	}

	// TODO (elliotcourant) Add a timeout here so that if the startup process for the connection
	//  takes too long we can simply fail early.
StartupLoop:
	for {
		response, err := connection.Receive()
		if err != nil {
			return nil, err
		}

		switch msg := response.(type) {
		case *pgproto.ErrorResponse:
			// TODO (elliotcourant) add proper error message handling/parsing.
			return nil, fmt.Errorf("pg: %s", msg.Message)
		case *pgproto.Authentication:
			switch msg.Type {
			// The authentication is okay, but we still need to wait for a ready for query
			// message to come through.
			case pgproto.AuthTypeOk:
			case pgproto.AuthTypeCleartextPassword:
				// If we get a cleartext password prompt then send the password we have as is and
				// do not encrypt it.
				if err := connection.Send(&pgproto.PasswordMessage{
					Password: dataNode.Password,
				}); err != nil {
					return nil, err
				}
			case pgproto.AuthTypeMD5Password:
				// If we get an MD5 password prompt, then we need to use the salt to encrypt the
				// password and send it back to the database.

				md5s := func(s string) string {
					h := md5.Sum([]byte(s))
					return hex.EncodeToString(h[:])
				}
				secret := "md5" + md5s(md5s(dataNode.Password+dataNode.Username)+string(msg.Salt[:]))
				if err := connection.Send(&pgproto.PasswordMessage{
					Password: secret,
				}); err != nil {
					return nil, err
				}
			default:
				// If we end up receiving another type of authentication that we cannot handle
				// then log a message about it.
				// TODO (elliotcourant) Add information to the log about the connection.
				timber.Warningf("received unexpected authentication type: %d", msg.Type)
			}
		case *pgproto.ReadyForQuery:
			// Once we receive a ready for query then we can use the connection.
			break StartupLoop
		default:
			// If we receive any other messages then we can simply discard them, for the most part
			// there aren't any other messages we would want to handle here.
		}
	}

	// If the startup loop was successful then we can return the connection here.
	return connection, nil
}

func (c *dataNodeShardConnection) DataNodeID() uint64 {
	panic("implement me")
}

func (c *dataNodeShardConnection) ShardID() uint64 {
	panic("implement me")
}

func (c *dataNodeShardConnection) DataNodeShardID() uint64 {
	panic("implement me")
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
