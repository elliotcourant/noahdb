package core

import (
	"fmt"
	"github.com/elliotcourant/noahdb/pkg/pgproto"
	"github.com/readystock/golog"
	"net"
	"sync"
)

type poolContext struct {
	*base
}

type poolItem struct {
	id   uint64
	pool sync.Pool
}

func (p *poolItem) addConnection(frontend *pgproto.Frontend) {
	p.pool.Put(frontend)
}

func (p *poolItem) GetConnection() PoolConnection {
	item := p.pool.Get().(*pgproto.Frontend)
	return &frontendConnection{
		item,
		p,
	}
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
	f.pool.pool.Put(f)
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

	}()
}

func (ctx *poolContext) GetConnectionForDataNodeShard(id uint64) (PoolConnection, error) {
	ctx.poolSync.Lock()
	defer ctx.poolSync.Unlock()
	pItem, ok := ctx.pool[id]
	if !ok {
		pItem = &poolItem{
			id:   id,
			pool: sync.Pool{},
		}
		ctx.pool[id] = pItem

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
				"user": "postgres",
			},
		}); err != nil {
			golog.Errorf("could not send startup message to data node [%d]: %s", dataNode.DataNodeID, err.Error())
			return nil, err
		}
		_, _ = frontend.Receive()

		ctx.pool[id].addConnection(frontend)
	}

	return pItem.GetConnection(), nil
}
