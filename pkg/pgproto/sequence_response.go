package pgproto

import (
	"bytes"
	"encoding/binary"
	"github.com/elliotcourant/noahdb/pkg/pgio"
)

type SequenceResponse struct {
	CurrentValue       uint64
	LastPartitionIndex uint64
	MaxPartitionIndex  uint64
	Partitions         uint64
}

func (sequence *SequenceResponse) Decode(src []byte) error {
	*sequence = SequenceResponse{}
	buf := bytes.NewBuffer(src)

	sequence.CurrentValue = binary.BigEndian.Uint64(buf.Next(8))
	sequence.LastPartitionIndex = binary.BigEndian.Uint64(buf.Next(8))
	sequence.MaxPartitionIndex = binary.BigEndian.Uint64(buf.Next(8))
	sequence.Partitions = binary.BigEndian.Uint64(buf.Next(8))

	return nil
}

func (sequence *SequenceResponse) Encode(dst []byte) []byte {
	dst = append(dst, RpcSequenceResponse)
	sp := len(dst)
	dst = pgio.AppendInt32(dst, -1)

	dst = pgio.AppendUint64(dst, sequence.CurrentValue)
	dst = pgio.AppendUint64(dst, sequence.LastPartitionIndex)
	dst = pgio.AppendUint64(dst, sequence.MaxPartitionIndex)
	dst = pgio.AppendUint64(dst, sequence.Partitions)

	pgio.SetInt32(dst[sp:], int32(len(dst[sp:])))

	return dst
}
