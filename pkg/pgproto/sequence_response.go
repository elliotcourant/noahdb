package pgproto

import (
	"github.com/elliotcourant/buffers"
)

type SequenceChunkResponse struct {
	Start  uint64
	End    uint64
	Offset uint64
	Count  uint64
}

func (item *SequenceChunkResponse) Decode(src []byte) error {
	*item = SequenceChunkResponse{}
	buf := buffers.NewBytesReader(src)

	item.Start = buf.NextUint64()
	item.End = buf.NextUint64()
	item.Offset = buf.NextUint64()
	item.Count = buf.NextUint64()

	return nil
}

func (item *SequenceChunkResponse) Encode(dst []byte) []byte {
	buf := buffers.NewBytesBuffer()
	buf.AppendByte(RpcSequenceResponse)
	buf.Append(item.EncodeBody()...)
	dst = append(dst, buf.Bytes()...)
	return dst
}

func (item *SequenceChunkResponse) EncodeBody() []byte {
	buf := buffers.NewBytesBuffer()
	buf.AppendUint64(item.Start)
	buf.AppendUint64(item.End)
	buf.AppendUint64(item.Offset)
	buf.AppendUint64(item.Count)
	return buf.Bytes()
}

func (SequenceChunkResponse) Backend() {}

type Sequence struct {
	CurrentValue       uint64
	LastPartitionIndex uint64
	MaxPartitionIndex  uint64
	Partitions         uint64
}

func (item *Sequence) Decode(src []byte) error {
	*item = Sequence{}
	buf := buffers.NewBytesReader(src)

	item.CurrentValue = buf.NextUint64()
	item.LastPartitionIndex = buf.NextUint64()
	item.MaxPartitionIndex = buf.NextUint64()
	item.Partitions = buf.NextUint64()

	return nil
}

// func (sequence *Sequence) Encode(dst []byte) []byte {
// 	dst = append(dst, RpcSequenceResponse)
// 	sp := len(dst)
// 	dst = pgio.AppendInt32(dst, -1)
//
// 	dst = pgio.AppendUint64(dst, sequence.CurrentValue)
// 	dst = pgio.AppendUint64(dst, sequence.LastPartitionIndex)
// 	dst = pgio.AppendUint64(dst, sequence.MaxPartitionIndex)
// 	dst = pgio.AppendUint64(dst, sequence.Partitions)
//
// 	pgio.SetInt32(dst[sp:], int32(len(dst[sp:])))
//
// 	return dst
// }

func (item *Sequence) EncodeBody() []byte {
	buf := buffers.NewBytesBuffer()
	buf.AppendUint64(item.CurrentValue)
	buf.AppendUint64(item.LastPartitionIndex)
	buf.AppendUint64(item.MaxPartitionIndex)
	buf.AppendUint64(item.Partitions)
	return buf.Bytes()
}
