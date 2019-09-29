package core

import (
	"github.com/elliotcourant/noahdb/pkg/core/static"
	"github.com/elliotcourant/noahdb/pkg/frunk"
	"github.com/elliotcourant/timber"
	"net"
	"sync"
)

type base struct {
	db *frunk.Store

	trans    TransportWrapper
	poolSync sync.RWMutex
	pool     map[uint64]*poolItem

	joinCluster func() error
}

func (ctx *base) State() frunk.ClusterState {
	if ctx.db == nil {
		return frunk.Unknown
	}
	return ctx.db.State()
}

func (ctx *base) LeaderID() (string, string, error) {
	addr := ctx.db.LeaderAddr()
	id, err := ctx.db.LeaderID()
	return addr, id, err
}

func (ctx *base) Neighbors() ([]*frunk.Server, error) {
	return ctx.db.Nodes()
}

func (ctx *base) JoinCluster() error {
	if ctx.joinCluster != nil {
		return ctx.joinCluster()
	}
	return nil
}

func (ctx *base) Join(id, addr string) error {
	return ctx.db.Join(id, addr, map[string]string{})
}

// Addr returns the address of the current node.
func (ctx *base) Addr() net.Addr {
	return ctx.trans.Addr()
}

// CoordinatorID returns the unique ID for this noahdb coordinator within the cluster.
func (ctx *base) CoordinatorID() uint64 {
	return uint64(1)
	// return ctx.db
}

// Close shuts down the colony.
func (ctx *base) Close() {
	ctx.db.Close(false)
}

// IsLeader returns true if the current coordinator is the leader of the cluster.
func (ctx *base) IsLeader() bool {
	return ctx.db.IsLeader()
}

// Setup initializes the internal store with any necessary data needed.
func (ctx *base) Setup(config ColonyConfig) {
	if !ctx.IsLeader() {
		return
	}

	if ctx.isSetup() {
		timber.Verbosef("internal database appears to be setup already")
		return
	}

	if _, err := ctx.db.ExecEx("PRAGMA journal_mode=WAL;", false); err != nil {
		panic(err)
	}

	setupScript, err := static.GetEmbeddedFile("/00_internal_sql.sql")
	if err != nil {
		panic(err)
	}

	_, err = ctx.db.Exec(string(setupScript))
	if err != nil {
		panic(err)
	}

	ctx.setupPostgresSystem()
}

func (ctx *base) Query(query string) (*frunk.QueryResponse, error) {
	return ctx.db.Query(query)
}

func (ctx *base) isSetup() bool {
	re, err := ctx.db.Query("SELECT data_node_id FROM data_nodes LIMIT 1;")
	// If the error is nil then that means the table exists and the cluster has been
	// setup. If the error is not nil then the cluster needs to be setup again.
	return err == nil && re.Rows[0].Error == ""
}
