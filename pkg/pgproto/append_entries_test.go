package pgproto

import (
	"encoding/hex"
	"fmt"
	"github.com/hashicorp/raft"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAppendEntriesRequest(t *testing.T) {
	t.Run("encode and decode", func(t *testing.T) {
		appendEntry := AppendEntriesRequest{
			AppendEntriesRequest: raft.AppendEntriesRequest{
				RPCHeader: raft.RPCHeader{
					ProtocolVersion: raft.ProtocolVersionMax,
				},
				Term:         5,
				Leader:       []byte("127.0.0.1:5432"),
				PrevLogEntry: 1003,
				PrevLogTerm:  5,
				Entries: []*raft.Log{
					{
						Index: 21342,
						Term:  5,
						Type:  raft.LogCommand,
						Data:  []byte("do all the things"),
					},
					{
						Index: 21343,
						Term:  5,
						Type:  raft.LogCommand,
						Data:  []byte("do all the things some more"),
					},
				},
				LeaderCommitIndex: 21341,
			},
		}
		encoded := appendEntry.Encode(nil)
		fmt.Println(hex.Dump(encoded))
		decodeEntry := AppendEntriesRequest{}
		err := decodeEntry.Decode(encoded[5:])
		assert.NoError(t, err)
		assert.Equal(t, appendEntry, decodeEntry)
	})

	t.Run("encode with nil and decode", func(t *testing.T) {
		appendEntry := AppendEntriesRequest{
			AppendEntriesRequest: raft.AppendEntriesRequest{
				RPCHeader: raft.RPCHeader{
					ProtocolVersion: raft.ProtocolVersionMax,
				},
				Term:         5,
				Leader:       []byte("127.0.0.1:5432"),
				PrevLogEntry: 1003,
				PrevLogTerm:  5,
				Entries: []*raft.Log{
					nil,
					{
						Index: 21343,
						Term:  5,
						Type:  raft.LogCommand,
						Data:  []byte("do all the things some more"),
					},
				},
				LeaderCommitIndex: 21341,
			},
		}
		encoded := appendEntry.Encode(nil)
		fmt.Println(hex.Dump(encoded))
		decodeEntry := AppendEntriesRequest{}
		err := decodeEntry.Decode(encoded[5:])
		assert.NoError(t, err)
		assert.Equal(t, appendEntry, decodeEntry)
	})
}
