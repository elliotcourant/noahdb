package sql

import (
	"github.com/elliotcourant/noahdb/pkg/ast"
	"github.com/elliotcourant/noahdb/pkg/pgproto"
	"github.com/elliotcourant/noahdb/pkg/pgwirebase"
	"github.com/elliotcourant/noahdb/pkg/types"
	"github.com/elliotcourant/noahdb/pkg/util/queryutil"
	"github.com/elliotcourant/noahdb/pkg/util/stmtbuf"
)

type sessionContext interface {
	Backend() *pgproto.Backend
	StatementBuffer() stmtbuf.StatementBuffer
}

type session struct {
	sessionContext

	preparedStatements map[string]preparedStatementEntry
	portals            map[string]portalEntry
}

func newSession(s sessionContext) *session {
	return &session{
		sessionContext:     s,
		preparedStatements: map[string]preparedStatementEntry{},
		portals:            map[string]portalEntry{},
	}
}

type preparedStatementEntry struct {
	*PreparedStatement
	portals map[string]struct{}
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
	Statement *ast.Stmt

	Types queryutil.PlaceholderTypes

	Columns []pgproto.FieldDescription

	InTypes []types.OID

	// TODO(andrei): The connExecutor doesn't use this. Delete it once the
	// Executor is gone.
	portalNames map[string]struct{}
}
