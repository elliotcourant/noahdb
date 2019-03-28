package types_test

import (
	"testing"
	"time"

	"github.com/elliotcourant/noahdb/pkg/types"
	"github.com/elliotcourant/noahdb/pkg/types/testutil"
)

func TestDaterangeTranscode(t *testing.T) {
	testutil.TestSuccessfulTranscodeEqFunc(t, "daterange", []interface{}{
		&types.Daterange{LowerType: types.Empty, UpperType: types.Empty, Status: types.Present},
		&types.Daterange{
			Lower:     types.Date{Time: time.Date(1990, 12, 31, 0, 0, 0, 0, time.UTC), Status: types.Present},
			Upper:     types.Date{Time: time.Date(2028, 1, 1, 0, 0, 0, 0, time.UTC), Status: types.Present},
			LowerType: types.Inclusive,
			UpperType: types.Exclusive,
			Status:    types.Present,
		},
		&types.Daterange{
			Lower:     types.Date{Time: time.Date(1800, 12, 31, 0, 0, 0, 0, time.UTC), Status: types.Present},
			Upper:     types.Date{Time: time.Date(2200, 1, 1, 0, 0, 0, 0, time.UTC), Status: types.Present},
			LowerType: types.Inclusive,
			UpperType: types.Exclusive,
			Status:    types.Present,
		},
		&types.Daterange{Status: types.Null},
	}, func(aa, bb interface{}) bool {
		a := aa.(types.Daterange)
		b := bb.(types.Daterange)

		return a.Status == b.Status &&
			a.Lower.Time.Equal(b.Lower.Time) &&
			a.Lower.Status == b.Lower.Status &&
			a.Lower.InfinityModifier == b.Lower.InfinityModifier &&
			a.Upper.Time.Equal(b.Upper.Time) &&
			a.Upper.Status == b.Upper.Status &&
			a.Upper.InfinityModifier == b.Upper.InfinityModifier
	})
}

func TestDaterangeNormalize(t *testing.T) {
	testutil.TestSuccessfulNormalizeEqFunc(t, []testutil.NormalizeTest{
		{
			SQL: "select daterange('2010-01-01', '2010-01-11', '(]')",
			Value: types.Daterange{
				Lower:     types.Date{Time: time.Date(2010, 1, 2, 0, 0, 0, 0, time.UTC), Status: types.Present},
				Upper:     types.Date{Time: time.Date(2010, 1, 12, 0, 0, 0, 0, time.UTC), Status: types.Present},
				LowerType: types.Inclusive,
				UpperType: types.Exclusive,
				Status:    types.Present,
			},
		},
	}, func(aa, bb interface{}) bool {
		a := aa.(types.Daterange)
		b := bb.(types.Daterange)

		return a.Status == b.Status &&
			a.Lower.Time.Equal(b.Lower.Time) &&
			a.Lower.Status == b.Lower.Status &&
			a.Lower.InfinityModifier == b.Lower.InfinityModifier &&
			a.Upper.Time.Equal(b.Upper.Time) &&
			a.Upper.Status == b.Upper.Status &&
			a.Upper.InfinityModifier == b.Upper.InfinityModifier
	})
}
