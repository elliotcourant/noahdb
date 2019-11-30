package engine

import (
	"github.com/elliotcourant/meles"
	"github.com/elliotcourant/mellivora"
)

var (
	_ Core        = &coreBase{}
	_ Transaction = &transactionBase{}
)

type (
	// Core is an interface around the base layer of noahdb.
	Core interface {
		Begin() (Transaction, error)
		Close() error
	}

	// Transaction is an interface that allows accessors to NoahDB's stateful models and allows
	// ways to make changes to them atomically.
	Transaction interface {
		Commit() error
		Rollback() error

		DataNodes() DataNodeContext
		DataNodeShards() DataNodeShardContext
	}

	coreBase struct {
		store *meles.Store
		db    *mellivora.Database
	}

	transactionBase struct {
		core *coreBase
		txn  *mellivora.Transaction
	}
)

// NewCore will create the Core interface from a meles store
// and a mellivora database.
func NewCore(store *meles.Store, db *mellivora.Database) Core {
	return &coreBase{
		store: store,
		db:    db,
	}
}

// Begin will start a new transaction.
func (c *coreBase) Begin() (Transaction, error) {
	txn, err := c.db.Begin()
	if err != nil {
		return nil, err
	}

	return &transactionBase{
		core: c,
		txn:  txn,
	}, nil
}

// Close will shut down this core and stop any pending actions.
func (c *coreBase) Close() error {
	return c.store.Stop()
}

// Commit will attempt to persist the changes made throughout the current transaction. If the commit
// is successful then the error will be nil. If it fails then an error will be returned.
func (t *transactionBase) Commit() error {
	return t.txn.Commit()
}

// Rollback will discard the changes in this transaction.
func (t *transactionBase) Rollback() error {
	return t.txn.Rollback()
}
