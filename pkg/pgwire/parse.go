package pgwire

import (
	"encoding/json"
	"fmt"
	"github.com/elliotcourant/noahdb/pkg/ast"
	"github.com/elliotcourant/noahdb/pkg/commands"
	"github.com/elliotcourant/noahdb/pkg/pgproto"
	"github.com/readystock/golog"
)

func (wire *wireServer) handleParse(parseMessage *pgproto.Parse) error {
	parseTree, err := ast.Parse(parseMessage.Query)
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
			Err: fmt.Errorf("could not handle statement [%s]", parseMessage.Query),
		})
	}

	j, _ := json.Marshal(parseTree)
	golog.Verbosef("received query: %s | %s", parseMessage.Query, string(j))
	if err := wire.StatementBuffer().Push(commands.PrepareStatement{
		Name:      parseMessage.Name,
		Statement: stmt,
	}); err != nil {
		return wire.StatementBuffer().Push(commands.SendError{
			Err: err,
		})
	}
	return wire.backend.Send(&pgproto.ParseComplete{})
}
