package pgproto

import (
	"bytes"
	"github.com/elliotcourant/noahdb/pkg/pgio"
)

// RpcNextSequence is the request used for nodes to request more sequence values.
type RpcNextSequence struct {
	SequenceName string
}

func (*RpcNextSequence) Frontend() {}

func (rpc *RpcNextSequence) Decode(src []byte) error {
	i := bytes.IndexByte(src, 0)
	if i != len(src)-1 {
		return &invalidMessageFormatErr{messageType: "Query"}
	}

	rpc.SequenceName = string(src[:i])

	return nil
}

func (rpc *RpcNextSequence) Encode(dst []byte) []byte {
	dst = append(dst, 's')
	dst = pgio.AppendInt32(dst, int32(4+len(rpc.SequenceName)+1))

	dst = append(dst, rpc.SequenceName...)
	dst = append(dst, 0)

	return dst
}
