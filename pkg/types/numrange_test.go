package types_test

import (
	"math/big"
	"testing"

	"github.com/elliotcourant/noahdb/pkg/types"
	"github.com/elliotcourant/noahdb/pkg/types/testutil"
)

func TestNumrangeTranscode(t *testing.T) {
	testutil.TestSuccessfulTranscode(t, "numrange", []interface{}{
		&types.Numrange{
			LowerType: types.Empty,
			UpperType: types.Empty,
			Status:    types.Present,
		},
		&types.Numrange{
			Lower:     types.Numeric{Int: big.NewInt(-543), Exp: 3, Status: types.Present},
			Upper:     types.Numeric{Int: big.NewInt(342), Exp: 1, Status: types.Present},
			LowerType: types.Inclusive,
			UpperType: types.Exclusive,
			Status:    types.Present,
		},
		&types.Numrange{
			Lower:     types.Numeric{Int: big.NewInt(-42), Exp: 1, Status: types.Present},
			Upper:     types.Numeric{Int: big.NewInt(-5), Exp: 0, Status: types.Present},
			LowerType: types.Inclusive,
			UpperType: types.Exclusive,
			Status:    types.Present,
		},
		&types.Numrange{
			Lower:     types.Numeric{Int: big.NewInt(-42), Exp: 1, Status: types.Present},
			LowerType: types.Inclusive,
			UpperType: types.Unbounded,
			Status:    types.Present,
		},
		&types.Numrange{
			Upper:     types.Numeric{Int: big.NewInt(-42), Exp: 1, Status: types.Present},
			LowerType: types.Unbounded,
			UpperType: types.Exclusive,
			Status:    types.Present,
		},
		&types.Numrange{Status: types.Null},
	})
}
