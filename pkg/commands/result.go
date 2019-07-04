package commands

import (
	"fmt"
	"github.com/elliotcourant/noahdb/pkg/ast"
	"github.com/elliotcourant/noahdb/pkg/pgproto"
	"github.com/elliotcourant/noahdb/pkg/util/stmtbuf"
)

type sessionContext interface {
	Backend() *pgproto.Backend
	StatementBuffer() stmtbuf.StatementBuffer
}

type completionMsgType int

const (
	_ completionMsgType = iota
	commandComplete
	bindComplete
	closeComplete
	parseComplete
	emptyQueryResponse
	readyForQuery
	flush
	// Some commands, like Describe, don't need a completion message.
	noCompletionMsg
)

type CommandResult struct {
	closed  bool
	session sessionContext
	err     error
	typ     completionMsgType
	tag     string

	noDataMessage bool
}

func CreateSyncCommandResult(session sessionContext) *CommandResult {
	result := NewCommandResult(session)
	result.typ = readyForQuery
	return result
}

func CreateExecuteCommandResult(session sessionContext, stmt ast.Stmt) *CommandResult {
	result := NewCommandResult(session)
	result.typ = commandComplete
	result.tag = stmt.StatementTag()
	return result
}

func CreatePreparedStatementResult(session sessionContext, stmt ast.Stmt) *CommandResult {
	result := NewCommandResult(session)
	result.typ = parseComplete
	result.tag = stmt.StatementTag()
	return result
}

func CreateDescribeStatementResult(session sessionContext) *CommandResult {
	result := NewCommandResult(session)
	result.typ = noCompletionMsg
	return result
}

func CreateErrorResult(session sessionContext, err error) *CommandResult {
	result := NewCommandResult(session)
	result.typ = noCompletionMsg
	result.err = err
	return result
}

func NewCommandResult(session sessionContext) *CommandResult {
	return &CommandResult{
		closed:  false,
		session: session,
	}
}

func (result *CommandResult) SetError(err error) {
	result.err = err
}

func (result *CommandResult) Err() error {
	return result.err
}

func (result *CommandResult) CloseWithErr(e error) error {
	if result.closed {
		return fmt.Errorf("command result is closed")
	}
	defer func() {
		result.closed = true
	}()
	return result.session.Backend().Send(&pgproto.ErrorResponse{
		Message: e.Error(),
	})
}

func (result *CommandResult) SetNoDataMessage(msg bool) {
	result.noDataMessage = msg
}

func (result *CommandResult) Close() error {
	if result.closed {
		return fmt.Errorf("command result is closed")
	}
	defer func() {
		result.closed = true
	}()

	// Send a completion message, specific to the type of result.
	switch result.typ {
	case commandComplete:
		tag := result.tag
		if err := result.session.Backend().Send(&pgproto.CommandComplete{
			CommandTag: tag,
		}); err != nil {
			panic(fmt.Sprintf("unexpected error from buffer: %s", err.Error()))
		}
		// panic("not handling command complete yet.")

		// tag := cookTag(
		// 	result.cmdCompleteTag, r.conn.writerState.tagBuf[:0], r.stmt, r.rowsAffected,
		// )
		// r.conn.bufferCommandComplete(tag)
	case parseComplete:
		if err := result.session.Backend().Send(&pgproto.ParseComplete{}); err != nil {
			panic(fmt.Sprintf("unexpected error from buffer: %s", err.Error()))
		}
	case bindComplete:
		if err := result.session.Backend().Send(&pgproto.BindComplete{}); err != nil {
			panic(fmt.Sprintf("unexpected error from buffer: %s", err.Error()))
		}
	case closeComplete:
		if err := result.session.Backend().Send(&pgproto.CloseComplete{}); err != nil {
			panic(fmt.Sprintf("unexpected error from buffer: %s", err.Error()))
		}
	case readyForQuery:
		if err := result.session.Backend().Send(&pgproto.ReadyForQuery{
			TxStatus: 'I',
		}); err != nil {
			panic(fmt.Sprintf("unexpected error from buffer: %s", err.Error()))
		}
	case emptyQueryResponse:
		if err := result.session.Backend().Send(&pgproto.EmptyQueryResponse{}); err != nil {
			panic(fmt.Sprintf("unexpected error from buffer: %s", err.Error()))
		}
	case flush:
		// // The error is saved on conn.err.
		// _ /* err */ = r.conn.Flush(r.pos)
	case noCompletionMsg:
		// Only for describe statements.
		if result.noDataMessage {
			if err := result.session.Backend().Send(&pgproto.NoData{}); err != nil {
				panic(fmt.Sprintf("unexpected error from buffer: %s", err.Error()))
			}
		}
	default:
		panic(fmt.Sprintf("unknown type: %v", result.typ))
	}
	return nil
}
