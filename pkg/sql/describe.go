package sql

import (
	"fmt"
	"github.com/elliotcourant/noahdb/pkg/ast"
	"github.com/elliotcourant/noahdb/pkg/commands"
	"github.com/elliotcourant/noahdb/pkg/pgwirebase"
)

func (s *session) ExecuteDescribe(describe commands.DescribeStatement, result *commands.CommandResult) error {
	ps, ok := s.preparedStatements[describe.Name]
	switch describe.Type {
	case pgwirebase.PrepareStatement:
		if !ok {
			return fmt.Errorf("unknown prepared statement %q", describe.Name)
		}

		if ps.Statement == nil || (*ps.Statement).StatementType() != ast.Rows {
			// The statement has no data to be returned.
			result.SetNoDataMessage(true)
		} else {

		}
	case pgwirebase.PreparePortal:
		if !ok {
			return fmt.Errorf("unknown portal %q", describe.Name)
		}

	default:
		return fmt.Errorf("unknown describe type: %s", describe.Type)
	}
	return nil
}
