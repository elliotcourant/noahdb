package pgproto

import (
	"encoding/hex"
	"fmt"
	"github.com/hashicorp/raft"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRequestVoteRequest(t *testing.T) {
	t.Run("encode and decode", func(t *testing.T) {
		requestVote := RequestVoteRequest{
			RequestVoteRequest: raft.RequestVoteRequest{
				RPCHeader: raft.RPCHeader{
					ProtocolVersion: raft.ProtocolVersionMax,
				},
				Term:         5,
				Candidate:    []byte("127.0.0.1:5432"),
				LastLogIndex: 1003,
				LastLogTerm:  5,
			},
		}
		encoded := requestVote.Encode(nil)
		fmt.Println(hex.Dump(encoded))
		decodeEntry := RequestVoteRequest{}
		err := decodeEntry.Decode(encoded[5:])
		assert.NoError(t, err)
		assert.Equal(t, requestVote, decodeEntry)
	})
}
