package pgproto

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/elliotcourant/noahdb/pkg/pgio"
	"github.com/hashicorp/raft"
)

type AppendEntriesResponse struct {
	raft.AppendEntriesResponse
	Error error
}

func (AppendEntriesResponse) Backend() {}

func (AppendEntriesResponse) Raft() {}

func (response *AppendEntriesResponse) Decode(src []byte) error {
	*response = AppendEntriesResponse{}
	buf := bytes.NewBuffer(src)
	response.ProtocolVersion = raft.ProtocolVersion(binary.BigEndian.Uint32(buf.Next(4)))
	response.Term = binary.BigEndian.Uint64(buf.Next(8))
	response.LastLog = binary.BigEndian.Uint64(buf.Next(8))
	response.Success = buf.Next(1)[0] == 1
	response.NoRetryBackoff = buf.Next(1)[0] == 1

	errSizeBytes := buf.Next(4)
	errSize := int(int32(binary.BigEndian.Uint32(errSizeBytes)))

	if errSize == -1 {
		response.Error = nil
	} else {
		response.Error = errors.New(string(buf.Next(errSize)))
	}

	return nil
}

func (response *AppendEntriesResponse) Encode(dst []byte) []byte {
	dst = append(dst, RaftAppendEntriesResponse)
	sp := len(dst)
	dst = pgio.AppendInt32(dst, -1)

	dst = pgio.AppendInt32(dst, int32(response.ProtocolVersion))
	dst = pgio.AppendUint64(dst, response.Term)
	dst = pgio.AppendUint64(dst, response.LastLog)
	dst = pgio.AppendBool(dst, response.Success)
	dst = pgio.AppendBool(dst, response.NoRetryBackoff)

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
