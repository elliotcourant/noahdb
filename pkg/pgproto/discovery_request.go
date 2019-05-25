package pgproto

import (
	"github.com/elliotcourant/noahdb/pkg/pgio"
)

type DiscoveryRequest struct{}

func (DiscoveryRequest) Frontend() {}

func (DiscoveryRequest) RpcFrontend() {}

func (discovery *DiscoveryRequest) Decode(src []byte) error {
	*discovery = DiscoveryRequest{}

	return nil
}

func (discovery *DiscoveryRequest) Encode(dst []byte) []byte {
	dst = append(dst, RpcDiscoveryRequest)
	sp := len(dst)

	pgio.SetInt32(dst[sp:], int32(len(dst[sp:])))

	return dst
}
