package pgwire

import (
	"github.com/elliotcourant/noahdb/pkg/commands"
	"github.com/elliotcourant/noahdb/pkg/pgproto"
	"github.com/elliotcourant/noahdb/pkg/pgwirebase"
)

func (wire *wireServer) handleDescribe(describeMessage *pgproto.Describe) error {
	return wire.stmtBuf.Push(commands.DescribeStatement{
		Name: describeMessage.Name,
		Type: pgwirebase.PrepareType(describeMessage.ObjectType),
	})
}
