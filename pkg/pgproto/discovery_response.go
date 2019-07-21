package pgproto

import (
	"bytes"
	"encoding/binary"
	"github.com/elliotcourant/noahdb/pkg/pgio"
)

type DiscoveryResponse struct {
	LeaderAddr string
}

func (DiscoveryResponse) Backend() {}

func (DiscoveryResponse) RpcBackend() {}

func (discovery *DiscoveryResponse) Decode(src []byte) error {
	*discovery = DiscoveryResponse{}
	buf := bytes.NewBuffer(src)

	leaderLen := int(binary.BigEndian.Uint32(buf.Next(4)))
	discovery.LeaderAddr = string(buf.Next(leaderLen))

	return nil
}

func (discovery *DiscoveryResponse) Encode(dst []byte) []byte {
	dst = append(dst, RpcDiscoveryRequest)
	sp := len(dst)
	dst = pgio.AppendInt32(dst, -1)

	dst = pgio.AppendInt32(dst, int32(len(discovery.LeaderAddr)))
	dst = append(dst, discovery.LeaderAddr...)

	pgio.SetInt32(dst[sp:], int32(len(dst[sp:])))

	return dst
}
