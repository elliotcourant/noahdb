package stmtbuf

import (
	"fmt"
	"github.com/elliotcourant/noahdb/pkg/util/syncutil"
	"sync"
)

type Command interface {
}

type StatementBuffer interface {
	Push(Command) error
}

// CmdPos represents the index of a command relative to the start of a
// connection. The first command received on a connection has position 0.
type cmdPos int64

// StmtBuf maintains a list of commands that a SQL client has sent for execution
// over a network connection. The commands are SQL queries to be executed,
// statements to be prepared, etc. At any point in time the buffer contains
// outstanding commands that have yet to be executed, and it can also contain
// some history of commands that we might want to retry - in the case of a
// retriable error, we'd like to retry all the commands pertaining to the
// current SQL transaction.
//
// The buffer is supposed to be used by one reader and one writer. The writer
// adds commands to the buffer using Push(). The reader reads one command at a
// time using curCmd(). The consumer is then supposed to create command results
// (the buffer is not involved in this).
// The buffer internally maintains a cursor representing the reader's position.
// The reader has to manually move the cursor using advanceOne(),
// seekToNextBatch() and rewind().
// In practice, the writer is a module responsible for communicating with a SQL
// client (i.e. pgwire.conn) and the reader is a connExecutor.
//
// The StmtBuf supports grouping commands into "batches" delimited by sync
// commands. A reader can then at any time chose to skip over commands from the
// current batch. This is used to implement Postgres error semantics: when an
// error happens during processing of a command, some future commands might need
// to be skipped. Batches correspond either to multiple queries received in a
// single query string (when the SQL client sends a semicolon-separated list of
// queries as part of the "simple" protocol), or to different commands pipelined
// by the cliend, separated from "sync" messages.
//
// push() can be called concurrently with curCmd().
//
// The connExecutor will use the buffer to maintain a window around the
// command it is currently executing. It will maintain enough history for
// executing commands again in case of an automatic retry. The connExecutor is
// in charge of trimming completed commands from the buffer when it's done with
// them.
type stmtBuf struct {
	syncutil.Mutex

	// closed, if set, means that the writer has closed the buffer. See Close().
	closed bool

	// cond is signaled when new commands are pushed.
	cond *sync.Cond

	// readerBlocked is set while the reader is blocked on waiting for a command
	// to be pushed into the buffer.
	readerBlocked bool

	// data contains the elements of the buffer.
	data []Command

	// startPos indicates the index of the first command currently in data
	// relative to the start of the connection.
	startPos cmdPos
	// curPos is the current position of the cursor going through the commands.
	// At any time, curPos indicates the position of the command to be returned
	// by curCmd().
	curPos cmdPos
	// lastPos indicates the position of the last command that was pushed into
	// the buffer.
	lastPos cmdPos
}

// Push adds a Command to the end of the buffer. If a curCmd() call was blocked
// waiting for this command to arrive, it will be woken up.
//
// An error is returned if the buffer has been closed.
func (buf *stmtBuf) Push(cmd Command) error { // ctx context.Context,
	buf.Lock()
	defer buf.Unlock()
	if buf.closed {
		return fmt.Errorf("buffer is closed")
	}
	buf.data = append(buf.data, cmd)
	buf.lastPos++

	buf.cond.Signal()
	return nil
}
