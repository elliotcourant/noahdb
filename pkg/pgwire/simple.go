package pgwire

import (
	"fmt"
	"github.com/elliotcourant/noahdb/pkg/ast"
	"github.com/elliotcourant/noahdb/pkg/commands"
	"github.com/elliotcourant/noahdb/pkg/pgproto"
)

func (wire *wireServer) handleSimpleQuery(parseMessage *pgproto.Query) error {
	parseTree, err := ast.Parse(parseMessage.String)
	if err != nil {
		return err
	}

	if len(parseTree.Statements) == 0 {
		// no statements
	}

	// Convert the actual query sent to the statement interface
	stmt, ok := parseTree.Statements[0].(ast.RawStmt).Stmt.(ast.Stmt)
	if !ok {
		return wire.StatementBuffer().Push(commands.SendError{
			Err: fmt.Errorf("could not handle statement [%s]", parseMessage.String),
		})
	}

	return wire.StatementBuffer().Push(commands.ExecuteStatement{
		Statement: stmt,
	})
}
