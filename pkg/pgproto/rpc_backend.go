package pgproto

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/jackc/pgx/chunkreader"
	"io"
)

type RpcBackend struct {
	cr *chunkreader.ChunkReader
	b  *bytes.Buffer
	w  io.Writer

	join JoinRequest

	bodyLen    int
	msgType    byte
	partialMsg bool
}

func NewRpcBackend(r io.Reader, w io.Writer) (*RpcBackend, error) {
	cr := chunkreader.NewChunkReader(r)
	return &RpcBackend{cr: cr, w: w, b: bytes.NewBuffer(make([]byte, 0))}, nil
}

func (b *RpcBackend) Send(msg BackendMessage) error {
	_, err := b.w.Write(msg.Encode(nil))
	return err
}

func (b *RpcBackend) Flush() error {
	_, err := b.b.WriteTo(b.w)
	return err
}

func (b *RpcBackend) Receive() (RpcFrontendMessage, error) {
	if !b.partialMsg {
		header, err := b.cr.Next(5)
		if err != nil {
			return nil, err
		}

		b.msgType = header[0]
		b.bodyLen = int(binary.BigEndian.Uint32(header[1:])) - 4
		b.partialMsg = true
	}

	var msg RpcFrontendMessage
	switch b.msgType {
	case RpcJoinRequest:
		msg = &b.join
	default:
		return nil, fmt.Errorf("unknown message type: %c", b.msgType)
	}

	msgBody, err := b.cr.Next(b.bodyLen)
	if err != nil {
		return nil, err
	}

	b.partialMsg = false

	err = msg.Decode(msgBody)
	return msg, err
}
