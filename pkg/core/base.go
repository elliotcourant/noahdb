package core

import (
	"github.com/elliotcourant/noahdb/pkg/core/static"
	"github.com/elliotcourant/noahdb/pkg/store"
	"github.com/readystock/golog"
)

type base struct {
	db *store.Store
}

// CoordinatorID returns the unique ID for this noahdb coordinator within the cluster.
func (ctx *base) CoordinatorID() uint64 {
	return ctx.db.NodeID()
}

// Close shuts down the colony.
func (ctx *base) Close() {
	ctx.db.Close()
}

// IsLeader returns true if the current coordinator is the leader of the cluster.
func (ctx *base) IsLeader() bool {
	return ctx.db.IsLeader()
}

// Setup initializes the internal store with any necessary data needed.
func (ctx *base) Setup() {
	if !ctx.IsLeader() {
		return
	}

	if ctx.isSetup() {
		golog.Verbosef("internal database appears to be setup already")
		return
	}

	setupScript, err := static.GetEmbeddedFile("/internal_sql.sql")
	if err != nil {
		panic(err)
	}

	_, err = ctx.db.Exec(string(setupScript))
	if err != nil {
		panic(err)
	}
}

func (ctx *base) isSetup() bool {
	_, err := ctx.db.Query("SELECT data_node_id FROM data_nodes LIMIT 1;")
	// If the error is nil then that means the table exists and the cluster has been
	// setup. If the error is not nil then the cluster needs to be setup again.
	return err == nil
}
