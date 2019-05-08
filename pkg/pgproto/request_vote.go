package pgproto

import (
	"bytes"
	"encoding/binary"
	"github.com/elliotcourant/noahdb/pkg/pgio"
	"github.com/hashicorp/raft"
)

type RequestVoteRequest struct {
	raft.RequestVoteRequest
}

func (RequestVoteRequest) Frontend() {}

func (RequestVoteRequest) RaftFrontend() {}

func (requestVote *RequestVoteRequest) Decode(src []byte) error {
	*requestVote = RequestVoteRequest{}
	buf := bytes.NewBuffer(src)
	requestVote.ProtocolVersion = raft.ProtocolVersion(binary.BigEndian.Uint32(buf.Next(4)))
	requestVote.Term = binary.BigEndian.Uint64(buf.Next(8))

	candidateLen := int(binary.BigEndian.Uint32(buf.Next(4)))
	requestVote.Candidate = buf.Next(candidateLen)

	requestVote.LastLogIndex = binary.BigEndian.Uint64(buf.Next(8))
	requestVote.LastLogTerm = binary.BigEndian.Uint64(buf.Next(8))

	return nil
}

func (requestVote *RequestVoteRequest) Encode(dst []byte) []byte {
	dst = append(dst, RaftRequestVoteRequest)
	sp := len(dst)
	dst = pgio.AppendInt32(dst, -1)

	dst = pgio.AppendInt32(dst, int32(requestVote.ProtocolVersion))
	dst = pgio.AppendUint64(dst, requestVote.Term)

	dst = pgio.AppendInt32(dst, int32(len(requestVote.Candidate)))
	dst = append(dst, requestVote.Candidate...)

	dst = pgio.AppendUint64(dst, requestVote.LastLogIndex)
	dst = pgio.AppendUint64(dst, requestVote.LastLogTerm)

	pgio.SetInt32(dst[sp:], int32(len(dst[sp:])))

	return dst
}
