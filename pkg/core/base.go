package core

import (
	"github.com/elliotcourant/noahdb/pkg/store"
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

func (ctx *base) Setup() {

}
