package pgproto

import (
	"github.com/elliotcourant/buffers"
)

type SequenceRequest struct {
	SequenceName string
}

func (SequenceRequest) Frontend() {}

func (SequenceRequest) RpcFrontend() {}

func (item *SequenceRequest) Encode(dst []byte) []byte {
	buf := buffers.NewBytesBuffer()
	buf.AppendByte(RpcSequenceRequest)
	buf.Append(item.EncodeBody()...)
	dst = append(dst, buf.Bytes()...)
	return dst
}

func (item *SequenceRequest) EncodeBody() []byte {
	buf := buffers.NewBytesBuffer()
	buf.AppendString(item.SequenceName)
	return buf.Bytes()
}

func (item *SequenceRequest) Decode(src []byte) error {
	*item = SequenceRequest{}
	buf := buffers.NewBytesReader(src)
	item.SequenceName = buf.NextString()
	// buf := bytes.NewBuffer(src)
	return nil
}
