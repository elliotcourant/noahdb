package hive

import (
	"github.com/elliotcourant/noahdb/pkg/core"
)

type Core interface {
	Base
	Begin() (Transaction, error)
}

type Base interface {
	NodeID() string
}

type Transaction interface {
	Base
	Commit() error
	Rollback() error

	DataNodes() DataNodeContext
}

type DataNodeContext interface {
	GetDataNodes() []core.DataNode
}
