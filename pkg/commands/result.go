package commands

import (
	"fmt"
	"github.com/elliotcourant/noahdb/pkg/ast"
	"github.com/elliotcourant/noahdb/pkg/pgproto"
)

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
	backend *pgproto.Backend
	err     error
	typ     completionMsgType
	tag     string

	noDataMessage bool
}

func CreateSyncCommandResult(backend *pgproto.Backend) *CommandResult {
	result := NewCommandResult(backend)
	result.typ = readyForQuery
	return result
}

func CreateExecuteCommandResult(backend *pgproto.Backend, stmt ast.Stmt) *CommandResult {
	result := NewCommandResult(backend)
	result.typ = commandComplete
	result.tag = stmt.StatementTag()
	return result
}

func CreatePreparedStatementResult(backend *pgproto.Backend, stmt ast.Stmt) *CommandResult {
	result := NewCommandResult(backend)
	result.typ = parseComplete
	result.tag = stmt.StatementTag()
	return result
}

func CreateDescribeStatementResult(backend *pgproto.Backend) *CommandResult {
	result := NewCommandResult(backend)
	result.typ = noCompletionMsg
	return result
}

func CreateBindStatementResult(backend *pgproto.Backend) *CommandResult {
	result := NewCommandResult(backend)
	result.typ = bindComplete
	return result
}

func CreateErrorResult(backend *pgproto.Backend, err error) *CommandResult {
	result := NewCommandResult(backend)
	result.typ = noCompletionMsg
	result.err = err
	return result
}

func NewCommandResult(backend *pgproto.Backend) *CommandResult {
	return &CommandResult{
		closed:  false,
		backend: backend,
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
	return result.backend.Send(&pgproto.ErrorResponse{
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
		if err := result.backend.Send(&pgproto.CommandComplete{
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
		if err := result.backend.Send(&pgproto.ParseComplete{}); err != nil {
			panic(fmt.Sprintf("unexpected error from buffer: %s", err.Error()))
		}
	case bindComplete:
		if err := result.backend.Send(&pgproto.BindComplete{}); err != nil {
			panic(fmt.Sprintf("unexpected error from buffer: %s", err.Error()))
		}
	case closeComplete:
		if err := result.backend.Send(&pgproto.CloseComplete{}); err != nil {
			panic(fmt.Sprintf("unexpected error from buffer: %s", err.Error()))
		}
	case readyForQuery:
		if err := result.backend.Send(&pgproto.ReadyForQuery{
			TxStatus: 'I',
		}); err != nil {
			panic(fmt.Sprintf("unexpected error from buffer: %s", err.Error()))
		}
	case emptyQueryResponse:
		if err := result.backend.Send(&pgproto.EmptyQueryResponse{}); err != nil {
			panic(fmt.Sprintf("unexpected error from buffer: %s", err.Error()))
		}
	case flush:
		// // The error is saved on conn.err.
		// _ /* err */ = r.conn.Flush(r.pos)
	case noCompletionMsg:
		// Only for describe statements.
		if result.noDataMessage {
			if err := result.backend.Send(&pgproto.NoData{}); err != nil {
				panic(fmt.Sprintf("unexpected error from buffer: %s", err.Error()))
			}
		}
	default:
		panic(fmt.Sprintf("unknown type: %v", result.typ))
	}
	return nil
}
