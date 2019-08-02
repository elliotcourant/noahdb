package sql

import (
	"github.com/elliotcourant/noahdb/pkg/ast"
	"github.com/elliotcourant/noahdb/pkg/pgwirebase"
	"github.com/elliotcourant/noahdb/pkg/util/queryutil"
)

type execResult interface {
	SetError(error)
	Err() error
}

func (s *session) executeStatement(
	stmt ast.Stmt,
	result execResult,
	placeholders queryutil.QueryArguments,
	outFormats []pgwirebase.FormatCode) error {
	result.SetError(s.stageQueryToResult(stmt, placeholders, outFormats))
	return nil
}
