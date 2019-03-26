package store

import (
	"encoding/hex"
	"fmt"
	"testing"
)

func TestUint64ToBytes(t *testing.T) {
	uints := []uint64{1, 2, 3, 15831904231, 35183541}
	for _, u := range uints {
		bytes := Uint64ToBytes(u)
		unt := BytesToUint64(bytes)
		fmt.Printf("Val: %d\nBytes: %sBack: %d\n\n", u, hex.Dump(bytes), unt)
		if unt != u {
			t.Errorf("could not convert bytes for %d back to uint64", u)
			t.Fail()
			return
		}
	}
}
