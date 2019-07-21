package pgproto

import (
	"encoding/hex"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCommandRequest(t *testing.T) {
	t.Run("encode and decode", func(t *testing.T) {
		item := CommandRequest{
			CommandType: RpcCommandType_Query,
			Queries: []string{
				"SELECT 1",
				"DELETE THINGS",
			},
			KeyValueSets: []KeyValue{
				{
					Key:   []byte("thing"),
					Value: []byte("value"),
				},
				{
					Key:   []byte("12412"),
					Value: []byte("sdgasgads"),
				},
			},
			KeyValueDeletes: [][]byte{
				[]byte("test"),
				[]byte("asnjkldgjkas"),
			},
		}
		encoded := item.Encode(nil)
		fmt.Println(hex.Dump(encoded))
		decodeEntry := CommandRequest{}
		err := decodeEntry.Decode(encoded[5:])
		assert.NoError(t, err)
		assert.Equal(t, item, decodeEntry)
	})
}
