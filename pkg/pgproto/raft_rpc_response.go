package pgproto

import (
	"encoding/json"
	"github.com/elliotcourant/noahdb/pkg/pgio"
	"github.com/hashicorp/raft"
)

type RaftRpcResponse struct {
	raft.RPCResponse
}

func (RaftRpcResponse) Backend() {}

func (RaftRpcResponse) RaftBackend() {}

func (response *RaftRpcResponse) Decode(src []byte) error {
	return nil
}

func (response *RaftRpcResponse) Encode(dst []byte) []byte {
	dst = append(dst, RaftInstallSnapshotRequest)
	sp := len(dst)
	dst = pgio.AppendInt32(dst, -1)

	if response.Error == nil {
		dst = pgio.AppendInt32(dst, -1)
	} else {
		errorBytes := []byte(response.Error.Error())
		dst = pgio.AppendInt32(dst, int32(len(errorBytes)))
		dst = append(dst, errorBytes...)
	}

	if response.Response != nil {
		dst = pgio.AppendInt32(dst, -1)
	} else {
		respBytes, _ := json.Marshal(response.Response)
		dst = pgio.AppendInt32(dst, int32(len(respBytes)))
		dst = append(dst, respBytes...)
	}

	pgio.SetInt32(dst[sp:], int32(len(dst[sp:])))
	return dst
}
