package sql

import (
	"github.com/elliotcourant/noahdb/pkg/commands"
	"github.com/elliotcourant/noahdb/pkg/pgproto"
)

func (s *session) ExecuteStatement(statement commands.ExecuteStatement, result *commands.CommandResult) error {
	err := s.stageQueryToResult(statement.Statement, result)
	if err != nil {
		return err
	}
	return s.Backend().Send(&pgproto.CommandComplete{
		CommandTag: statement.Statement.StatementTag(),
	})
}
