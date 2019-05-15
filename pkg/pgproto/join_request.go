package pgproto

import (
	"bytes"
	"encoding/binary"
	"github.com/elliotcourant/noahdb/pkg/pgio"
)

type JoinRequest struct {
	NodeID  string
	Address string
}

func (JoinRequest) Frontend() {}

func (join *JoinRequest) Decode(src []byte) error {
	*join = JoinRequest{}
	buf := bytes.NewBuffer(src)
	nodeIdLen := int(binary.BigEndian.Uint32(buf.Next(4)))
	join.NodeID = string(buf.Next(nodeIdLen))

	addressLen := int(binary.BigEndian.Uint32(buf.Next(4)))
	join.Address = string(buf.Next(addressLen))

	return nil
}

func (join *JoinRequest) Encode(dst []byte) []byte {
	dst = append(dst, RpcJoinRequest)
	sp := len(dst)
	dst = pgio.AppendInt32(dst, -1)

	dst = pgio.AppendInt32(dst, int32(len(join.NodeID)))
	dst = append(dst, join.NodeID...)

	dst = pgio.AppendInt32(dst, int32(len(join.Address)))
	dst = append(dst, join.Address...)

	pgio.SetInt32(dst[sp:], int32(len(dst[sp:])))

	return dst
}
