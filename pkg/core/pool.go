package core

import (
	"fmt"
	"github.com/elliotcourant/noahdb/pkg/pgproto"
	"github.com/readystock/golog"
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

type frontendInterface interface {
	Send(pgproto.FrontendMessage) error
	Receive() (pgproto.BackendMessage, error)
}

type frontendConnection struct {
	*pgproto.Frontend

	pool *poolItem
}

func (f *frontendConnection) Release() {
	if f.Frontend == nil {
		return
	}
	golog.Verbosef("releasing connection from data node shard [%d], pool size: %d", f.pool.id, len(f.pool.pool))
	f.pool.releaseConnection(f)
}

type PoolConnection interface {
	frontendInterface
	Release()
}

type PoolContext interface {
	GetConnectionForDataNodeShard(id uint64) (PoolConnection, error)
}

func (ctx *base) Pool() PoolContext {
	return &poolContext{
		ctx,
	}
}

func (ctx *base) StartPool() {
	go func() {
		for {
			time.Sleep(30 * time.Second)
			if !ctx.IsLeader() {
				continue
			}

			// dataNodeShards, err := ctx.DataNodes().GetDataNodes()
		}
	}()
}

func (ctx *poolContext) GetConnectionForDataNodeShard(id uint64) (PoolConnection, error) {
	ctx.poolSync.Lock()
	defer ctx.poolSync.Unlock()
	pItem, ok := ctx.pool[id]
	if !ok {
		golog.Tracef("data node shard [%d] is not in pool, creating connection", id)
		pItem = &poolItem{
			id:    id,
			mutex: sync.Mutex{},
			pool:  make([]*frontendConnection, 0),
		}
		ctx.pool[id] = pItem

		conn, err := ctx.NewConnection(id)
		if err != nil {
			golog.Errorf("could not create connection to data node shard [%d]: %v", id, err)
			return nil, err
		}

		ctx.pool[id].addConnection(conn)
	}
	poolConn := pItem.GetConnection()
	if poolConn == nil {
		conn, err := ctx.NewConnection(id)
		if err != nil {
			golog.Errorf("could not create connection to data node shard [%d]: %v", id, err)
			return nil, err
		}
		return &frontendConnection{
			Frontend: conn,
			pool:     pItem,
		}, nil
	}
	return poolConn, nil
}

func (ctx *poolContext) NewConnection(id uint64) (*pgproto.Frontend, error) {
	dataNode, err := ctx.DataNodes().GetDataNodeForDataNodeShard(id)
	if err != nil {
		return nil, err
	}

	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", dataNode.Address, dataNode.Port))
	if err != nil {
		golog.Errorf("could not resolve address for data node [%d]: %s", dataNode.DataNodeID, err.Error())
		return nil, err
	}

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		golog.Errorf("could not connect to data node [%d]: %s", dataNode.DataNodeID, err.Error())
		return nil, err
	}

	frontend, err := pgproto.NewFrontend(conn, conn)
	if err != nil {
		golog.Errorf("could not setup frontend for data node [%d]: %s", dataNode.DataNodeID, err.Error())
		return nil, err
	}

	if err := frontend.Send(&pgproto.StartupMessage{
		ProtocolVersion: pgproto.ProtocolVersionNumber,
		Parameters: map[string]string{
			"user":     "postgres",
			"database": fmt.Sprintf("noahdb_%d", id),
		},
	}); err != nil {
		golog.Errorf("could not send startup message to data node [%d]: %s", dataNode.DataNodeID, err.Error())
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
					panic("authentication is not implemented")
				}
			case *pgproto.ParameterStatus:
			case *pgproto.ParameterDescription:
			case *pgproto.BackendKeyData:
			case *pgproto.ReadyForQuery:
				return nil // We are good to go, exit the loop
			case *pgproto.ErrorResponse:
				return fmt.Errorf("from backend: %s", msg.Message)
			default:
				golog.Warnf("unexpected message from backend %T", msg)
			}
		}
	}(); err != nil {
		return nil, err
	}

	return frontend, nil
}
