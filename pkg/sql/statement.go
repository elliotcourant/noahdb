package sql

import (
	"github.com/elliotcourant/noahdb/pkg/commands"
)

type execResult interface {
	SetError(error)
	Err() error
}

func (s *session) ExecuteStatement(statement commands.ExecuteStatement, result execResult) error {
	result.SetError(s.stageQueryToResult(statement.Statement))
	return nil
}
