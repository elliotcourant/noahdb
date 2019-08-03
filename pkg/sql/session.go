package sql

import (
	"github.com/elliotcourant/noahdb/pkg/ast"
	"github.com/elliotcourant/noahdb/pkg/core"
	"github.com/elliotcourant/noahdb/pkg/pgproto"
	"github.com/elliotcourant/noahdb/pkg/pgwirebase"
	"github.com/elliotcourant/noahdb/pkg/types"
	"github.com/elliotcourant/noahdb/pkg/util/queryutil"
	"github.com/elliotcourant/noahdb/pkg/util/stmtbuf"
	"github.com/elliotcourant/timber"
	"sync"
)

type QueryMode int

const (
	QueryModeStandard QueryMode = 0
	QueryModeExtended           = 1
)

type TransactionState int

const (
	TransactionState_None   TransactionState = 0
	TransactionState_Active                  = 1
)

type sessionContext interface {
	Backend() *pgproto.Backend
	Colony() core.Colony
	StatementBuffer() stmtbuf.StatementBuffer
}

type session struct {
	sessionContext

	preparedStatements   map[string]preparedStatementEntry
	portals              map[string]portalEntry
	log                  timber.Logger
	queryMode            QueryMode
	queryModeSync        sync.RWMutex
	transactionState     TransactionState
	transactionStateSync sync.RWMutex

	pool     map[uint64]core.PoolConnection
	poolSync sync.Mutex
}

func (s *session) SetQueryMode(mode QueryMode) {
	s.queryModeSync.Lock()
	defer s.queryModeSync.Unlock()
	s.queryMode = mode
}

func (s *session) GetQueryMode() QueryMode {
	s.queryModeSync.RLock()
	defer s.queryModeSync.RUnlock()
	return s.queryMode
}

func (s *session) SetTransactionState(state TransactionState) {
	s.transactionStateSync.Lock()
	defer s.transactionStateSync.Unlock()
	s.transactionState = state
}

func (s *session) GetTransactionState() TransactionState {
	s.transactionStateSync.RLock()
	defer s.transactionStateSync.RUnlock()
	return s.transactionState
}

func (s *session) GetConnectionForDataNodeShard(id uint64) (core.PoolConnection, error) {
	s.transactionStateSync.Lock()
	defer s.transactionStateSync.Unlock()
	if pool, ok := s.pool[id]; ok {
		return pool, nil
	}
	pc, err := s.Colony().Pool().GetConnectionForDataNodeShard(id)
	if err != nil {
		return nil, err
	}
	s.pool[id] = pc

	return pc, nil
}

func (s *session) ReleaseConnectionForDataNodeShard(conn core.PoolConnection) {
	s.transactionStateSync.Lock()
	defer s.transactionStateSync.Unlock()
	if _, ok := s.pool[conn.ID()]; ok {
		delete(s.pool, conn.ID())
	}
	conn.Release()
}

func newSession(s sessionContext, log timber.Logger) *session {
	return &session{
		sessionContext:     s,
		preparedStatements: map[string]preparedStatementEntry{},
		portals:            map[string]portalEntry{},
		log:                log,
		pool:               map[uint64]core.PoolConnection{},
	}
}

type preparedStatementEntry struct {
	*PreparedStatement
}

// PreparedPortal is a PreparedStatement that has been bound with query arguments.
type PreparedPortal struct {
	Stmt  *PreparedStatement
	Qargs queryutil.QueryArguments

	// OutFormats contains the requested formats for the output columns.
	OutFormats []pgwirebase.FormatCode
}

// PreparedPortal is a PreparedStatement that has been bound with query arguments.
type portalEntry struct {
	*PreparedPortal
	psName string
}

// PreparedStatement is a SQL statement that has been parsed and the types
// of arguments and results have been determined.
type PreparedStatement struct {
	// Str is the statement string prior to parsing, used to generate
	// error messages. This may be used in
	// the future to present a contextual error message based on location
	// information.
	Str string

	// TypeHints contains the types of the placeholders set by the client. It
	// dictates how input parameters for those placeholders will be parsed. If a
	// placeholder has no type hint, it will be populated during type checking.
	TypeHints queryutil.PlaceholderTypes

	// Statement is the parse tree from pg_query.
	// This is used later to modify the query on the fly.
	Statement ast.Stmt

	Types queryutil.PlaceholderTypes

	Columns []pgproto.FieldDescription

	InferredTypes []types.Type
}
