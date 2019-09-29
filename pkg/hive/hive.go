package hive

import (
	"github.com/elliotcourant/meles"
	"github.com/elliotcourant/timber"
	"net"
)

type Core interface {
	Base
	Start() error
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

func NewHive(
	listener net.Listener, logger timber.Logger, options meles.Options,
) (Core, error) {
	db, err := meles.NewStore(listener, logger, options)
	if err != nil {
		return nil, err
	}
	return &hive{
		db: db,
	}, nil
}

type hive struct {
	db *meles.Store
}

func (h *hive) Start() error {
	return h.db.Start()
}

func (h *hive) NodeID() string {
	return ""
}

func (h *hive) Begin() (Transaction, error) {
	txn, err := h.db.Begin()
	if err != nil {
		return nil, err
	}
	return &hiveTransaction{
		txn:  txn,
		Base: h,
	}, nil
}

type hiveTransaction struct {
	txn *meles.Transaction
	Base
}

func (h *hiveTransaction) Commit() error {
	return h.txn.Commit()
}

func (h *hiveTransaction) Rollback() error {
	return h.txn.Rollback()
}
