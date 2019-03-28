package types_test

import (
	"testing"

	"github.com/elliotcourant/noahdb/pkg/types"
	"github.com/elliotcourant/noahdb/pkg/types/testutil"
)

func TestBitTranscode(t *testing.T) {
	testutil.TestSuccessfulTranscode(t, "bit(40)", []interface{}{
		&types.Varbit{Bytes: []byte{0, 0, 0, 0, 0}, Len: 40, Status: types.Present},
		&types.Varbit{Bytes: []byte{0, 1, 128, 254, 255}, Len: 40, Status: types.Present},
		&types.Varbit{Status: types.Null},
	})
}

func TestBitNormalize(t *testing.T) {
	testutil.TestSuccessfulNormalize(t, []testutil.NormalizeTest{
		{
			SQL:   "select B'111111111'",
			Value: &types.Bit{Bytes: []byte{255, 128}, Len: 9, Status: types.Present},
		},
	})
}
