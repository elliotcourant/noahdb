package pgwire

import (
	"encoding/json"
	"fmt"
	"github.com/elliotcourant/noahdb/pkg/ast"
	"github.com/elliotcourant/noahdb/pkg/commands"
	"github.com/elliotcourant/noahdb/pkg/pgproto"
	"github.com/elliotcourant/noahdb/pkg/types"
	"github.com/elliotcourant/noahdb/pkg/util/queryutil"
)

func (wire *Server) handleParse(parseMessage *pgproto.Parse) error {
	parseTree, err := ast.Parse(parseMessage.Query)
	if err != nil {
		return err
	}

	if len(parseTree.Statements) == 0 {
		// no statements
	} else if len(parseTree.Statements) > 1 {
		return fmt.Errorf("cannot have more than 1 statement per message in extended query")
	}

	rawTypeHints := make([]types.OID, len(parseMessage.ParameterOIDs))
	for i, raw := range parseMessage.ParameterOIDs {
		rawTypeHints[i] = types.OID(raw)
	}

	// Convert the actual query sent to the statement interface
	stmt, ok := parseTree.Statements[0].(ast.RawStmt).Stmt.(ast.Stmt)
	if !ok {
		return fmt.Errorf("could not handle statement [%s]", parseMessage.Query)
	}

	j, _ := json.Marshal(parseTree)
	wire.log.Verbosef("received query: %s | %s", parseMessage.Query, string(j))

	placeholders := queryutil.GetArguments(stmt)

	if len(rawTypeHints) > len(placeholders) {
		return fmt.Errorf("received too many type hints: %d vs %d placeholders in query",
			len(rawTypeHints), len(placeholders))
	}

	var sqlTypeHints queryutil.PlaceholderTypes
	if len(rawTypeHints) > 0 {
		// Prepare the mapping of SQL placeholder names to types. Pre-populate it with
		// the type hints received from the client, if any.
		sqlTypeHints = make(queryutil.PlaceholderTypes, len(placeholders))
		for i, t := range rawTypeHints {
			if t == 0 {
				continue
			}
			v, ok := types.GetTypeByOid(t)
			if !ok {
				return fmt.Errorf("unknown oid type: %v", t)
			}
			sqlTypeHints[i] = v
		}
	}

	return wire.StatementBuffer().Push(commands.PrepareStatement{
		Name:         parseMessage.Name,
		Statement:    stmt,
		TypeHints:    sqlTypeHints,
		RawTypeHints: rawTypeHints,
	})
}
