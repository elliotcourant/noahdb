package pgwire

import (
	"encoding/json"
	"github.com/elliotcourant/noahdb/pkg/ast"
	"github.com/elliotcourant/noahdb/pkg/pgwire/pgproto"
	"github.com/readystock/golog"
)

func (wire *wireServer) handleParse(parseMessage *pgproto.Parse) error {
	parseTree, err := ast.Parse(parseMessage.Query)
	if err != nil {
		return err
	}
	j, _ := json.Marshal(parseTree)
	golog.Verbosef("received query: %s | %s", parseMessage.Query, string(j))
	return wire.backend.Send(&pgproto.ParseComplete{})
}
