package pgproto

import (
	"encoding/binary"
	"fmt"
	"github.com/jackc/pgx/chunkreader"
	"io"
)

type RaftBackend struct {
	cr *chunkreader.ChunkReader
	w  io.Writer

	appendEntries AppendEntries

	bodyLen    int
	msgType    byte
	partialMsg bool
}

func NewRaftBackend(r io.Reader, w io.Writer) (*RaftBackend, error) {
	cr := chunkreader.NewChunkReader(r)
	return &RaftBackend{cr: cr, w: w}, nil
}

func (b *RaftBackend) Send(msg RaftBackendMessage) error {
	_, err := b.w.Write(msg.Encode(nil))
	return err
}

func (b *RaftBackend) Receive() (RaftFrontendMessage, error) {
	if !b.partialMsg {
		header, err := b.cr.Next(5)
		if err != nil {
			return nil, err
		}

		b.msgType = header[0]
		b.bodyLen = int(binary.BigEndian.Uint32(header[1:])) - 4
		b.partialMsg = true
	}

	var msg RaftFrontendMessage
	switch b.msgType {
	case RaftAppendEntriesRequest:
		msg = &b.appendEntries
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
