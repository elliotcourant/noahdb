package pgproto

import (
	"encoding/hex"
	"fmt"
	"github.com/hashicorp/raft"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAppendEntriesResponse(t *testing.T) {
	t.Run("encode and decode", func(t *testing.T) {
		appendEntry := AppendEntriesResponse{
			AppendEntriesResponse: raft.AppendEntriesResponse{
				RPCHeader: raft.RPCHeader{
					ProtocolVersion: raft.ProtocolVersionMax,
				},
				Term:           5,
				LastLog:        123,
				Success:        true,
				NoRetryBackoff: true,
			},
		}
		encoded := appendEntry.Encode(nil)
		fmt.Println(hex.Dump(encoded))
		decodeEntry := AppendEntriesResponse{}
		err := decodeEntry.Decode(encoded[5:])
		assert.NoError(t, err)
		assert.Equal(t, appendEntry, decodeEntry)
	})

	t.Run("encode and decode with error", func(t *testing.T) {
		appendEntry := AppendEntriesResponse{
			AppendEntriesResponse: raft.AppendEntriesResponse{
				RPCHeader: raft.RPCHeader{
					ProtocolVersion: raft.ProtocolVersionMax,
				},
				Term:           5,
				LastLog:        123,
				Success:        false,
				NoRetryBackoff: true,
			},
			Error: fmt.Errorf("test error"),
		}
		encoded := appendEntry.Encode(nil)
		fmt.Println(hex.Dump(encoded))
		decodeEntry := AppendEntriesResponse{}
		err := decodeEntry.Decode(encoded[5:])
		assert.NoError(t, err)
		assert.Equal(t, appendEntry, decodeEntry)
	})
}
