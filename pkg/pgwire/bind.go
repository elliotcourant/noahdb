package pgwire

import (
	"github.com/elliotcourant/noahdb/pkg/commands"
	"github.com/elliotcourant/noahdb/pkg/pgproto"
	"github.com/elliotcourant/noahdb/pkg/pgwirebase"
	"math"
)

func (wire *Server) handleBind(bindMessage *pgproto.Bind) error {
	qArgFormatCodes := make([]pgwirebase.FormatCode, int(math.Max(float64(len(bindMessage.ParameterFormatCodes)), 1)))
	switch len(bindMessage.ParameterFormatCodes) {
	case 0:
		// No format codes means all arguments are passed as text.
		qArgFormatCodes[0] = pgwirebase.FormatText
	default:
		for i := range qArgFormatCodes {
			qArgFormatCodes[i] = pgwirebase.FormatCode(bindMessage.ParameterFormatCodes[i])
		}
	}

	outFormats := make([]pgwirebase.FormatCode, int(math.Max(float64(len(bindMessage.ResultFormatCodes)), 1)))
	switch len(bindMessage.ResultFormatCodes) {
	case 0:
		outFormats[0] = pgwirebase.FormatText
	default:
		for i := range qArgFormatCodes {
			outFormats[i] = pgwirebase.FormatCode(bindMessage.ResultFormatCodes[i])
		}
	}

	return wire.stmtBuf.Push(commands.BindStatement{
		PreparedStatementName: bindMessage.PreparedStatement,
		PortalName:            bindMessage.DestinationPortal,
		Args:                  bindMessage.Parameters,
		ArgFormatCodes:        qArgFormatCodes,
		OutFormats:            outFormats,
	})
}
