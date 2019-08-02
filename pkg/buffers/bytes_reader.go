package buffers

import (
	"encoding/binary"
)

type BytesReader interface {
	NextByte() byte
	NextBytes() []byte
	NextString() string
	NextInt8() int8
	NextInt16() int16
	NextInt32() int32
	NextInt64() int64
	NextUint8() uint8
	NextUint16() uint16
	NextUint32() uint32
	NextUint64() uint64
	NextBool() bool
}

func NewBytesReader(src []byte) BytesReader {
	return &bytesReader{
		data:   src,
		offset: 0,
	}
}

type bytesReader struct {
	data   []byte
	offset uint32
}

func (b *bytesReader) NextByte() byte {
	i := b.data[b.offset]
	b.offset++
	return i
}

func (b *bytesReader) NextBytes() []byte {
	length := b.NextUint32()
	i := b.data[b.offset : b.offset+length]
	b.offset += length
	return i
}

func (b *bytesReader) NextString() string {
	return string(b.NextBytes())
}

func (b *bytesReader) NextInt8() int8 {
	return int8(b.NextUint8())
}

func (b *bytesReader) NextInt16() int16 {
	return int16(b.NextUint16())
}

func (b *bytesReader) NextInt32() int32 {
	return int32(b.NextUint32())
}

func (b *bytesReader) NextInt64() int64 {
	return int64(b.NextUint64())
}

func (b *bytesReader) NextUint8() uint8 {
	return b.NextByte()
}

func (b *bytesReader) NextBool() bool {
	t := b.NextByte()
	return t == 1
}

func (b *bytesReader) NextUint16() uint16 {
	i := b.data[b.offset : b.offset+Uint16Size]
	b.offset += Uint16Size
	return binary.BigEndian.Uint16(i)
}

func (b *bytesReader) NextUint32() uint32 {
	i := b.data[b.offset : b.offset+Uint32Size]
	b.offset += Uint32Size
	return binary.BigEndian.Uint32(i)
}

func (b *bytesReader) NextUint64() uint64 {
	i := b.data[b.offset : b.offset+Uint64Size]
	b.offset += Uint64Size
	return binary.BigEndian.Uint64(i)
}
