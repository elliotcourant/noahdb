package pgwire

import (
	"github.com/elliotcourant/noahdb/pkg/commands"
	"github.com/elliotcourant/noahdb/pkg/pgproto"
)

func (wire *Server) handleExecute(executeMessage *pgproto.Parse) error {
	return wire.stmtBuf.Push(commands.ExecutePortal{})
}
