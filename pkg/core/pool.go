package core

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/elliotcourant/noahdb/pkg/pgproto"
	"github.com/elliotcourant/timber"
	"net"
	"sync"
	"time"
)

type poolContext struct {
	*base
}

type poolItem struct {
	id    uint64
	mutex sync.Mutex
	pool  []*frontendConnection
}

func (p *poolItem) addConnection(frontend *pgproto.Frontend) {
	p.releaseConnection(&frontendConnection{
		Frontend: frontend,
		pool:     p,
	})
}

func (p *poolItem) releaseConnection(conn *frontendConnection) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.pool = append(p.pool, conn)
}

func (p *poolItem) GetConnection() PoolConnection {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if len(p.pool) == 0 {
		return nil
	}
	item := p.pool[0]
	p.pool = p.pool[1:]
	return item
}

func (p *poolItem) Size() int {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	return len(p.pool)
}

type frontendInterface interface {
	Send(pgproto.FrontendMessage) error
	Receive() (pgproto.BackendMessage, error)
	Close()
}

type frontendConnection struct {
	conn net.Conn
	*pgproto.Frontend

	pool *poolItem
}

func (f *frontendConnection) Release() {
	if f.Frontend == nil {
		return
	}
	timber.Verbosef("releasing connection from data node shard [%d], pool size: %d", f.pool.id, len(f.pool.pool))
	f.pool.releaseConnection(f)
}

func (f *frontendConnection) Close() {
	f.conn.Close()
	f.Frontend = nil
}

type PoolConnection interface {
	frontendInterface
	Release()
}

type PoolContext interface {
	StartPool()
	GetConnectionForDataNodeShard(id uint64) (PoolConnection, error)
}

func (ctx *base) Pool() PoolContext {
	return &poolContext{
		ctx,
	}
}

func (ctx *poolContext) StartPool() {
	desiredPoolSize := 5
	go func() {
		for {
			time.Sleep(30 * time.Second)
			if !ctx.IsLeader() {
				continue
			}

			dataNodeShards, err := ctx.Shards().GetDataNodeShards()
			if err != nil {
				timber.Errorf("could not retrieve data node shards for pool check: %v", err)
				continue
			}

			for _, dataNodeShard := range dataNodeShards {
				pool, err := ctx.getPoolForDataNodeShard(dataNodeShard.DataNodeShardID)
				if err != nil {
					timber.Errorf("could not retrieve pool for data node shard [%d], could not verify health: %v", dataNodeShard.DataNodeShardID, err)
					continue
				}

				size := pool.Size()

				if size == desiredPoolSize {
					timber.Verbosef("data node shard [%d] pool full, size: %d", dataNodeShard.DataNodeShardID, size)
					continue
				}

				if size < desiredPoolSize {
					// If the pool is not full then we should try to top it off.
					timber.Verbosef("data node shard [%d] pool not full, size: %d", dataNodeShard.DataNodeShardID, size)
					for i := size; i < desiredPoolSize; i++ {
						conn, err := ctx.newConnection(dataNodeShard.DataNodeShardID, pool)
						if err != nil {
							timber.Errorf("could not create connection to add to pool: %v", err)
							continue
						}
						// We've now created a new connection, release it to the pool for use.
						conn.Release()
					}
				} else {
					// If the pool is over flowing then grab some connections and throw them out.
					for i := size; i > desiredPoolSize; i-- {
						conn := pool.GetConnection()
						if conn != nil {
							conn.Close()
						}
					}
				}

				timber.Verbosef("data node shard [%d] new pool size: %d", dataNodeShard.DataNodeShardID, pool.Size())
			}
		}
	}()
}

func (ctx *poolContext) getPoolForDataNodeShard(id uint64) (*poolItem, error) {
	ctx.poolSync.RLock()
	pItem, ok := ctx.pool[id]
	ctx.poolSync.RUnlock()
	if !ok {
		timber.Tracef("data node shard [%d] is not in pool, creating connection", id)
		pItem = &poolItem{
			id:    id,
			mutex: sync.Mutex{},
			pool:  make([]*frontendConnection, 0),
		}
		ctx.poolSync.Lock()
		ctx.pool[id] = pItem
		ctx.poolSync.Unlock()

		// conn, err := ctx.NewConnection(id)
		// if err != nil {
		// 	timber.Errorf("could not create connection to data node shard [%d]: %v", id, err)
		// 	return nil, err
		// }
		// ctx.pool[id].addConnection(conn)
	}
	return pItem, nil
}

func (ctx *poolContext) GetConnectionForDataNodeShard(id uint64) (PoolConnection, error) {
	pItem, err := ctx.getPoolForDataNodeShard(id)
	if err != nil {
		return nil, err
	}
	poolConn := pItem.GetConnection()
	if poolConn == nil {
		return ctx.newConnection(id, pItem)
	}
	return poolConn, nil
}

func (ctx *poolContext) newConnection(id uint64, pool *poolItem) (*frontendConnection, error) {
	dataNode, err := ctx.DataNodes().GetDataNodeForDataNodeShard(id)
	if err != nil {
		return nil, err
	}

	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", dataNode.GetAddress(), dataNode.GetPort()))
	if err != nil {
		timber.Errorf("could not resolve address for data node [%d]: %s", dataNode.DataNodeID, err.Error())
		return nil, err
	}

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		timber.Errorf("could not connect to data node [%d]: %s", dataNode.DataNodeID, err.Error())
		return nil, err
	}

	frontend, err := pgproto.NewFrontend(conn, conn)
	if err != nil {
		timber.Errorf("could not setup frontend for data node [%d]: %s", dataNode.DataNodeID, err.Error())
		return nil, err
	}

	if err := frontend.Send(&pgproto.StartupMessage{
		ProtocolVersion: pgproto.ProtocolVersionNumber,
		Parameters: map[string]string{
			"user":     dataNode.GetUser(),
			"database": fmt.Sprintf("noahdb_%d", id),
		},
	}); err != nil {
		timber.Errorf("could not send startup message to data node [%d]: %s", dataNode.DataNodeID, err.Error())
		return nil, err
	}

	if err := func() error {
		for {
			response, err := frontend.Receive()
			if err != nil {
				return err
			}

			switch msg := response.(type) {
			case *pgproto.Authentication:
				if msg.Type != pgproto.AuthTypeOk {
					if msg.Type == pgproto.AuthTypeMD5Password {
						md5s := func(s string) string {
							h := md5.Sum([]byte(s))
							return hex.EncodeToString(h[:])
						}
						secret := "md5" + md5s(md5s(dataNode.GetPassword()+dataNode.GetUser())+string(msg.Salt[:]))
						if err := frontend.Send(&pgproto.PasswordMessage{
							Password: secret,
						}); err != nil {
							return fmt.Errorf("could not authenticate: %v", err)
						}
					} else {
						panic(fmt.Sprintf("cannot handle authentication type: %d", msg.Type))
					}
				}
			case *pgproto.ParameterStatus:
			case *pgproto.ParameterDescription:
			case *pgproto.BackendKeyData:
			case *pgproto.ReadyForQuery:
				return nil // We are good to go, exit the loop
			case *pgproto.ErrorResponse:
				return fmt.Errorf("from backend: %s", msg.Message)
			default:
				timber.Warningf("unexpected message from backend %T", msg)
			}
		}
	}(); err != nil {
		return nil, err
	}

	return &frontendConnection{
		Frontend: frontend,
		pool:     pool,
		conn:     conn,
	}, nil
}
