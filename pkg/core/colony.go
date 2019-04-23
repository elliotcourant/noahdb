package core

import (
	"github.com/elliotcourant/noahdb/pkg/store"
	"time"
)

// Colony is a wrapper for all of the core data that noahdb needs to operate.
type Colony interface {
	Shards() ShardContext
	Tenants() TenantContext
	DataNodes() DataNodeContext
	Tables() TableContext
	Schema() SchemaContext
	Users() UserContext
	// Shards()
	// Nodes()
	// Tenants()
	// Network()
	// Settings()
	// Pool()
	// Schema()
	// Sequences()

	CoordinatorID() uint64
	Close()
}

func NewColony(dataDirectory, listenAddress, joinAddress, postgresAddress string) (Colony, error) {
	db, err := store.CreateStore(dataDirectory, listenAddress, joinAddress)
	if err != nil {
		return nil, err
	}

	colony := &base{
		db: db,
	}

	time.Sleep(6 * time.Second)

	colony.Setup()

	return colony, nil
}
