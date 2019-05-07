package pgproto

import (
	"encoding/binary"
	"fmt"
	"github.com/elliotcourant/noahdb/pkg/pgio"
)

type RpcStartupMessage struct {
}

func (*RpcStartupMessage) Frontend() {}

func (dst *RpcStartupMessage) Decode(src []byte) error {
	if len(src) < 4 {
		return fmt.Errorf("startup message too short")
	}

	protocolVersion := binary.BigEndian.Uint32(src)
	if protocolVersion != RpcNumber {
		return fmt.Errorf("message is not a valid rpc startup message")
	}

	return nil
}

func (src *RpcStartupMessage) Encode(dst []byte) []byte {
	sp := len(dst)
	dst = pgio.AppendInt32(dst, -1)
	dst = pgio.AppendUint32(dst, RpcNumber)
	dst = append(dst, 0)
	pgio.SetInt32(dst[sp:], int32(len(dst[sp:])))
	return dst
}
