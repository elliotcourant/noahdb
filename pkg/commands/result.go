package commands

import (
	"fmt"
	"github.com/elliotcourant/noahdb/pkg/pgproto"
	"github.com/elliotcourant/noahdb/pkg/util/stmtbuf"
)

type sessionContext interface {
	Backend() *pgproto.Backend
	StatementBuffer() stmtbuf.StatementBuffer
}

type CommandResult struct {
	closed  bool
	session sessionContext
}

func NewCommandResult(session sessionContext) *CommandResult {
	return &CommandResult{
		closed:  false,
		session: session,
	}
}

func (result *CommandResult) CloseWithErr(err error) error {
	if result.closed {
		return fmt.Errorf("command result is closed")
	}
	defer func() {
		result.closed = true
	}()
	return result.session.Backend().Send(&pgproto.ErrorResponse{
		Message: err.Error(),
	})
}

func (result *CommandResult) Close() error {
	if result.closed {
		return fmt.Errorf("command result is closed")
	}
	defer func() {
		result.closed = true
	}()
	return result.session.Backend().Send(&pgproto.ReadyForQuery{
		TxStatus: 'I',
	})
}
