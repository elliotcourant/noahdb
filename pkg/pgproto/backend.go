package pgproto

import (
	"bytes"
	"encoding/binary"
	"github.com/elliotcourant/timber"
	"io"

	"github.com/jackc/pgx/chunkreader"
	"github.com/pkg/errors"
)

type Backend struct {
	cr *chunkreader.ChunkReader
	w  io.Writer
	b  *bytes.Buffer

	// Frontend message flyweights
	bind            Bind
	_close          Close
	describe        Describe
	execute         Execute
	flush           Flush
	parse           Parse
	passwordMessage PasswordMessage
	query           Query
	startupMessage  StartupMessage
	sync            Sync
	terminate       Terminate

	bodyLen    int
	msgType    byte
	partialMsg bool
}

func NewBackend(r io.Reader, w io.Writer) (*Backend, error) {
	cr := chunkreader.NewChunkReader(r)
	return &Backend{cr: cr, w: w, b: bytes.NewBuffer(make([]byte, 0))}, nil
}

func (b *Backend) Send(msg BackendMessage) error {
	timber.Verbosef("sending [%T]", msg)
	_, err := b.w.Write(msg.Encode(nil))
	return err
}

func (b *Backend) Flush() error {
	_, err := b.b.WriteTo(b.w)
	return err
}

func (b *Backend) ReceiveInitialMessage() (int, error) {
	buf, err := b.cr.Next(4)
	if err != nil {
		return 0, err
	}
	_ = int(binary.BigEndian.Uint32(buf) - 4)

	buf, err = b.cr.Next(4)
	if err != nil {
		return 0, err
	}
	return int(binary.BigEndian.Uint32(buf) - 4), nil
}

func (b *Backend) ReceiveStartupMessage() (*StartupMessage, error) {
	buf, err := b.cr.Next(4)
	if err != nil {
		return nil, err
	}
	msgSize := int(binary.BigEndian.Uint32(buf) - 4)

	buf, err = b.cr.Next(msgSize)
	if err != nil {
		return nil, err
	}

	err = b.startupMessage.Decode(buf)
	if err != nil {
		return nil, err
	}

	return &b.startupMessage, nil
}

func (b *Backend) Receive() (FrontendMessage, error) {
	if !b.partialMsg {
		header, err := b.cr.Next(5)
		if err != nil {
			return nil, err
		}

		b.msgType = header[0]
		b.bodyLen = int(binary.BigEndian.Uint32(header[1:])) - 4
		b.partialMsg = true
	}

	var msg FrontendMessage
	switch b.msgType {
	case PgBind:
		msg = &b.bind
	case PgClose:
		msg = &b._close
	case PgDescribe:
		msg = &b.describe
	case PgExecute:
		msg = &b.execute
	case PgFlush:
		msg = &b.flush
	case PgParse:
		msg = &b.parse
	case PgPasswordMessage:
		msg = &b.passwordMessage
	case PgQuery:
		msg = &b.query
	case PgSync:
		msg = &b.sync
	case PgTerminate:
		msg = &b.terminate
	default:
		return nil, errors.Errorf("unknown message type: %c", b.msgType)
	}

	msgBody, err := b.cr.Next(b.bodyLen)
	if err != nil {
		return nil, err
	}

	b.partialMsg = false

	err = msg.Decode(msgBody)
	return msg, err
}
