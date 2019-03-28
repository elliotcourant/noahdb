package stmtbuf

import (
	"fmt"
	"github.com/elliotcourant/noahdb/pkg/util/syncutil"
	"io"
	"sync"
)

type Command interface {
	Command()
}

type StatementBuffer interface {
	Push(Command) error
	AdvanceOne()
	CurrentCommand() (Command, CmdPos, error)
}

// CmdPos represents the index of a command relative to the start of a
// connection. The first command received on a connection has position 0.
type CmdPos int64

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
	startPos CmdPos
	// curPos is the current position of the cursor going through the commands.
	// At any time, curPos indicates the position of the command to be returned
	// by curCmd().
	curPos CmdPos
	// lastPos indicates the position of the last command that was pushed into
	// the buffer.
	lastPos CmdPos
}

func NewStatementBuffer() StatementBuffer {
	var buf stmtBuf
	buf.lastPos = -1
	buf.cond = sync.NewCond(&buf.Mutex)
	return &buf
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

// CurrentCommand returns the Command currently indicated by the cursor. Besides the
// Command itself, the command's position is also returned; the position can be
// used to later rewind() to this Command.
// If the cursor is positioned over an empty slot, the call blocks until the
// next Command is pushed into the buffer.
// If the buffer has previously been Close()d, or is closed while this is
// blocked, io.EOF is returned.
func (buf *stmtBuf) CurrentCommand() (Command, CmdPos, error) {
	buf.Lock()
	defer buf.Unlock()
	for {
		if buf.closed {
			return nil, 0, io.EOF
		}
		curPos := buf.curPos
		cmdIdx, err := buf.translatePosLocked(curPos)
		if err != nil {
			return nil, 0, err
		}
		if cmdIdx < len(buf.data) {
			return buf.data[cmdIdx], curPos, nil
		}
		if cmdIdx != len(buf.data) {
			return nil, 0, fmt.Errorf(
				"can only wait for next command; corrupt cursor: %d", curPos)
		}
		// Wait for the next Command to arrive to the buffer.
		buf.readerBlocked = true
		buf.cond.Wait()
		buf.readerBlocked = false
	}
}

// AdvanceOne advances the cursor one Command over. The command over which the
// cursor will be positioned when this returns may not be in the buffer yet.
func (buf *stmtBuf) AdvanceOne() {
	buf.Lock()
	buf.curPos++
	buf.Unlock()
}

// translatePosLocked translates an absolute position of a command (counting
// from the connection start) to the index of the respective command in the
// buffer (so, it returns an index relative to the start of the buffer).
//
// Attempting to translate a position that's below buf.startPos returns an
// error.
func (buf *stmtBuf) translatePosLocked(pos CmdPos) (int, error) {
	if pos < buf.startPos {
		return 0, fmt.Errorf(
			"position %d no longer in buffer (buffer starting at %d)",
			pos, buf.startPos)
	}
	return int(pos - buf.startPos), nil
}
