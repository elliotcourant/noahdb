package pgwire

import (
	"encoding/json"
	"github.com/elliotcourant/noahdb/pkg/ast"
	"github.com/elliotcourant/noahdb/pkg/pgwire/pgproto"
	"github.com/readystock/golog"
)

func (wire *wireServer) handleParse(name, query string, parameterOids []uint32) error {
	parseTree, err := ast.Parse(query)
	if err != nil {
		return err
	}
	j, _ := json.Marshal(parseTree)
	golog.Verbosef("received query: %s | %s", query, string(j))
	return wire.backend.Send(&pgproto.ParseComplete{})
}
