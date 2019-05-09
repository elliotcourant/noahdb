package pgproto

import (
	"encoding/binary"
	"fmt"
	"github.com/elliotcourant/noahdb/pkg/pgio"
)

type RaftStartupMessage struct {
}

func (RaftStartupMessage) Initial() {}

func (RaftStartupMessage) Frontend() {}

func (RaftStartupMessage) RaftFrontend() {}

func (dst *RaftStartupMessage) Decode(src []byte) error {
	if len(src) < 4 {
		return fmt.Errorf("startup message too short")
	}

	protocolVersion := binary.BigEndian.Uint32(src)
	if protocolVersion != RaftNumber {
		return fmt.Errorf("message is not a valid raft startup message")
	}

	return nil
}

func (src *RaftStartupMessage) Encode(dst []byte) []byte {
	sp := len(dst)
	dst = pgio.AppendInt32(dst, -1)
	dst = pgio.AppendUint32(dst, RaftNumber)
	dst = append(dst, 0)
	pgio.SetInt32(dst[sp:], int32(len(dst[sp:])))
	return dst
}
