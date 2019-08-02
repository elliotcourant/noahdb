package buffers

import (
	"encoding/binary"
)

const (
	Uint8Size  = 1
	Uint16Size = 2
	Uint32Size = 4
	Uint64Size = 8
)

type BytesBuffer interface {
	Append(bytes ...byte)
	AppendByte(byte byte)
	AppendString(str string)
	AppendUint8(item uint8)
	AppendUint16(item uint16)
	AppendUint32(item uint32)
	AppendUint64(item uint64)
	AppendInt16(item int16)
	AppendInt32(item int32)
	AppendInt64(item int64)
	AppendBool(item bool)
	AppendNil32()
	AppendBuffer(buffer BytesBuffer)

	Bytes() []byte
}

func NewBytesBuffer() BytesBuffer {
	return &bytesBuffer{
		buf: make([]byte, 0),
	}
}

type bytesBuffer struct {
	buf []byte
}

func (b *bytesBuffer) AppendBuffer(buffer BytesBuffer) {
	b.Append(buffer.Bytes()...)
}

func (b *bytesBuffer) AppendByte(byte byte) {
	b.AppendUint8(byte)
}

func (b *bytesBuffer) Append(bytes ...byte) {
	b.AppendUint32(uint32(len(bytes)))
	b.buf = append(b.buf, bytes...)
}

func (b *bytesBuffer) AppendString(str string) {
	b.Append([]byte(str)...)
}

func (b *bytesBuffer) AppendUint8(item uint8) {
	b.buf = append(b.buf, item)
}

func (b *bytesBuffer) AppendUint16(item uint16) {
	wp := len(b.buf)
	b.buf = append(b.buf, 0, 0)
	binary.BigEndian.PutUint16(b.buf[wp:], item)
}

func (b *bytesBuffer) AppendUint32(item uint32) {
	wp := len(b.buf)
	b.buf = append(b.buf, 0, 0, 0, 0)
	binary.BigEndian.PutUint32(b.buf[wp:], item)
}

func (b *bytesBuffer) AppendUint64(item uint64) {
	wp := len(b.buf)
	b.buf = append(b.buf, 0, 0, 0, 0, 0, 0, 0, 0)
	binary.BigEndian.PutUint64(b.buf[wp:], item)
}

func (b *bytesBuffer) AppendInt16(item int16) {
	b.AppendUint16(uint16(item))
}

func (b *bytesBuffer) AppendInt32(item int32) {
	b.AppendUint32(uint32(item))
}

func (b *bytesBuffer) AppendInt64(item int64) {
	b.AppendUint64(uint64(item))
}

func (b *bytesBuffer) AppendBool(item bool) {
	switch item {
	case true:
		b.AppendUint8(1)
	default:
		b.AppendUint8(0)
	}
}

func (b *bytesBuffer) AppendNil32() {
	b.AppendInt32(-1)
}

func (b *bytesBuffer) Bytes() []byte {
	return b.buf
}
