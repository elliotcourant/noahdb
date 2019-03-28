package types_test

import (
	"testing"

	"github.com/elliotcourant/noahdb/pkg/types"
	"github.com/elliotcourant/noahdb/pkg/types/testutil"
)

func TestInt4rangeTranscode(t *testing.T) {
	testutil.TestSuccessfulTranscode(t, "int4range", []interface{}{
		&types.Int4range{LowerType: types.Empty, UpperType: types.Empty, Status: types.Present},
		&types.Int4range{Lower: types.Int4{Int: 1, Status: types.Present}, Upper: types.Int4{Int: 10, Status: types.Present}, LowerType: types.Inclusive, UpperType: types.Exclusive, Status: types.Present},
		&types.Int4range{Lower: types.Int4{Int: -42, Status: types.Present}, Upper: types.Int4{Int: -5, Status: types.Present}, LowerType: types.Inclusive, UpperType: types.Exclusive, Status: types.Present},
		&types.Int4range{Lower: types.Int4{Int: 1, Status: types.Present}, LowerType: types.Inclusive, UpperType: types.Unbounded, Status: types.Present},
		&types.Int4range{Upper: types.Int4{Int: 1, Status: types.Present}, LowerType: types.Unbounded, UpperType: types.Exclusive, Status: types.Present},
		&types.Int4range{Status: types.Null},
	})
}

func TestInt4rangeNormalize(t *testing.T) {
	testutil.TestSuccessfulNormalize(t, []testutil.NormalizeTest{
		{
			SQL:   "select int4range(1, 10, '(]')",
			Value: types.Int4range{Lower: types.Int4{Int: 2, Status: types.Present}, Upper: types.Int4{Int: 11, Status: types.Present}, LowerType: types.Inclusive, UpperType: types.Exclusive, Status: types.Present},
		},
	})
}
