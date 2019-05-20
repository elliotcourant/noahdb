package pgproto

import (
	"encoding/binary"
	"fmt"
	"github.com/jackc/pgx/chunkreader"
	"io"
	"time"
)

type RaftWire struct {
	cr *chunkreader.ChunkReader
	w  io.Writer

	appendEntriesRequest    AppendEntriesRequest
	appendEntriesResponse   AppendEntriesResponse
	requestVoteRequest      RequestVoteRequest
	requestVoteResponse     RequestVoteResponse
	installSnapshotRequest  InstallSnapshotRequest
	installSnapshotResponse InstallSnapshotResponse

	errorResponse ErrorResponse

	bodyLen    int
	msgType    byte
	partialMsg bool
}

func NewRaftWire(r io.Reader, w io.Writer) (*RaftWire, error) {
	cr := chunkreader.NewChunkReader(r)
	return &RaftWire{cr: cr, w: w}, nil
}

func (b *RaftWire) Send(msg RaftMessage) error {
	time.Sleep(1 * time.Millisecond)
	// j, _ := json.Marshal(msg)
	// golog.Verbosef("sending %T message: %s", msg, string(j))
	_, err := b.w.Write(msg.Encode(nil))
	return err
}

func (b *RaftWire) Receive() (RaftMessage, error) {
	if !b.partialMsg {
		header, err := b.cr.Next(5)
		if err != nil {
			return nil, err
		}

		b.msgType = header[0]
		b.bodyLen = int(binary.BigEndian.Uint32(header[1:])) - 4
		b.partialMsg = true
	}

	var msg RaftMessage
	switch b.msgType {

	// Append entries
	case RaftAppendEntriesRequest:
		msg = &b.appendEntriesRequest
	case RaftAppendEntriesResponse:
		msg = &b.appendEntriesResponse

	// Request vote
	case RaftRequestVoteRequest:
		msg = &b.requestVoteRequest
	case RaftRequestVoteResponse:
		msg = &b.requestVoteResponse

	// Install snapshot
	case RaftInstallSnapshotRequest:
		msg = &b.installSnapshotRequest
	case RaftInstallSnapshotResponse:
		msg = &b.installSnapshotResponse

	case PgErrorResponse:
		msg = &b.errorResponse

	default:
		return nil, fmt.Errorf("unknown raft message type: %c", b.msgType)
	}

	msgBody, err := b.cr.Next(b.bodyLen)
	if err != nil {
		return nil, err
	}

	b.partialMsg = false

	err = msg.Decode(msgBody)
	time.Sleep(1 * time.Millisecond)
	// j, _ := json.Marshal(msg)
	// golog.Verbosef("received %T message: %s", msg, string(j))

	return msg, err
}
