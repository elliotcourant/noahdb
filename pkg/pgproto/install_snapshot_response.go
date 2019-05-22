package pgproto

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/elliotcourant/noahdb/pkg/pgio"
	"github.com/hashicorp/raft"
)

type InstallSnapshotResponse struct {
	raft.InstallSnapshotResponse
	Error error
}

func (InstallSnapshotResponse) Backend() {}

func (InstallSnapshotResponse) Raft() {}

func (response *InstallSnapshotResponse) Decode(src []byte) error {
	*response = InstallSnapshotResponse{}
	buf := bytes.NewBuffer(src)
	response.ProtocolVersion = raft.ProtocolVersion(binary.BigEndian.Uint32(buf.Next(4)))
	response.Term = binary.BigEndian.Uint64(buf.Next(8))
	response.Success = buf.Next(1)[0] == 1

	errSizeBytes := buf.Next(4)
	errSize := int(int32(binary.BigEndian.Uint32(errSizeBytes)))

	if errSize == -1 {
		response.Error = nil
	} else {
		response.Error = errors.New(string(buf.Next(errSize)))
	}

	return nil
}

func (response *InstallSnapshotResponse) Encode(dst []byte) []byte {
	dst = append(dst, RaftInstallSnapshotResponse)
	sp := len(dst)
	dst = pgio.AppendInt32(dst, -1)

	dst = pgio.AppendInt32(dst, int32(response.ProtocolVersion))
	dst = pgio.AppendUint64(dst, response.Term)
	dst = pgio.AppendBool(dst, response.Success)

	if response.Error == nil {
		dst = pgio.AppendInt32(dst, -1)
	} else {
		err := []byte(response.Error.Error())
		dst = pgio.AppendInt32(dst, int32(len(err)))
		dst = append(dst, err...)
	}

	pgio.SetInt32(dst[sp:], int32(len(dst[sp:])))

	return dst
}
