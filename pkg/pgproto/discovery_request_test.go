package pgproto

import (
	"encoding/hex"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDiscoveryRequest(t *testing.T) {
	t.Run("encode and decode", func(t *testing.T) {
		discover := DiscoveryRequest{}
		encoded := discover.Encode(nil)
		fmt.Println(hex.Dump(encoded))
		decodeEntry := DiscoveryRequest{}
		err := decodeEntry.Decode(encoded[5:])
		assert.NoError(t, err)
		assert.Equal(t, discover, decodeEntry)
	})
}
