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
		msg = &AppendEntriesRequest{}
	case RaftAppendEntriesResponse:
		msg = &AppendEntriesResponse{}
	// Request vote
	case RaftRequestVoteRequest:
		msg = &RequestVoteRequest{}
	case RaftRequestVoteResponse:
		msg = &RequestVoteResponse{}

	// Install snapshot
	case RaftInstallSnapshotRequest:
		msg = &InstallSnapshotRequest{}
	case RaftInstallSnapshotResponse:
		msg = &InstallSnapshotResponse{}

	case PgErrorResponse:
		msg = &ErrorResponse{}

	default:
		return nil, fmt.Errorf("unknown raft message type: %c", b.msgType)
	}

	msgBody, err := b.cr.Next(b.bodyLen)
	if err != nil {
		return nil, err
	}

	b.partialMsg = false

	err = msg.Decode(msgBody)

	return msg, err
}
