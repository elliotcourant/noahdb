package pgproto

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/elliotcourant/noahdb/pkg/pgio"
	"github.com/hashicorp/raft"
)

type RequestVoteResponse struct {
	raft.RequestVoteResponse
	Error error
}

func (RequestVoteResponse) Raft() {}

func (response *RequestVoteResponse) Decode(src []byte) error {
	*response = RequestVoteResponse{}
	buf := bytes.NewBuffer(src)
	response.ProtocolVersion = raft.ProtocolVersion(binary.BigEndian.Uint32(buf.Next(4)))
	response.Term = binary.BigEndian.Uint64(buf.Next(8))

	peersLength := int(binary.BigEndian.Uint32(buf.Next(4)))
	response.Peers = buf.Next(peersLength)

	response.Granted = buf.Next(1)[0] == 1

	errSizeBytes := buf.Next(4)
	errSize := int(int32(binary.BigEndian.Uint32(errSizeBytes)))

	if errSize == -1 {
		response.Error = nil
	} else {
		response.Error = errors.New(string(buf.Next(errSize)))
	}

	return nil
}

func (response *RequestVoteResponse) Encode(dst []byte) []byte {
	dst = append(dst, RaftRequestVoteResponse)
	sp := len(dst)
	dst = pgio.AppendInt32(dst, -1)

	dst = pgio.AppendInt32(dst, int32(response.ProtocolVersion))
	dst = pgio.AppendUint64(dst, response.Term)

	dst = pgio.AppendInt32(dst, int32(len(response.Peers)))
	dst = append(dst, response.Peers...)

	dst = pgio.AppendBool(dst, response.Granted)

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
