package sql

import (
	"github.com/elliotcourant/noahdb/pkg/ast"
)

type execResult interface {
	SetError(error)
	Err() error
}

func (s *session) executeStatement(
	stmt ast.Stmt,
	result execResult,
	placeholders interface{}) error {
	result.SetError(s.stageQueryToResult(stmt, placeholders))
	return nil
}
