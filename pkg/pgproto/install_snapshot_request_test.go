package pgproto

import (
	"encoding/hex"
	"fmt"
	"github.com/hashicorp/raft"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInstallSnapshotRequest(t *testing.T) {
	t.Run("encode and decode", func(t *testing.T) {
		item := InstallSnapshotRequest{
			InstallSnapshotRequest: raft.InstallSnapshotRequest{
				RPCHeader: raft.RPCHeader{
					ProtocolVersion: raft.ProtocolVersionMax,
				},
				SnapshotVersion:    raft.SnapshotVersionMax,
				Term:               5,
				Leader:             []byte("127.0.0.1:5432"),
				LastLogIndex:       1003,
				LastLogTerm:        5,
				Peers:              []byte("the stuff and the things"),
				Configuration:      []byte("other stuff and things"),
				ConfigurationIndex: 213,
				Size:               21495219431,
			},
		}
		encoded := item.Encode(nil)
		fmt.Println(hex.Dump(encoded))
		decodeEntry := InstallSnapshotRequest{}
		err := decodeEntry.Decode(encoded[5:])
		assert.NoError(t, err)
		assert.Equal(t, item, decodeEntry)
	})
}
