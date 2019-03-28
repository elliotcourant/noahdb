package types_test

import (
	"testing"

	"github.com/elliotcourant/noahdb/pkg/types"
	"github.com/elliotcourant/noahdb/pkg/types/testutil"
)

func TestIntervalTranscode(t *testing.T) {
	testutil.TestSuccessfulTranscode(t, "interval", []interface{}{
		&types.Interval{Microseconds: 1, Status: types.Present},
		&types.Interval{Microseconds: 1000000, Status: types.Present},
		&types.Interval{Microseconds: 1000001, Status: types.Present},
		&types.Interval{Microseconds: 123202800000000, Status: types.Present},
		&types.Interval{Days: 1, Status: types.Present},
		&types.Interval{Months: 1, Status: types.Present},
		&types.Interval{Months: 12, Status: types.Present},
		&types.Interval{Months: 13, Days: 15, Microseconds: 1000001, Status: types.Present},
		&types.Interval{Microseconds: -1, Status: types.Present},
		&types.Interval{Microseconds: -1000000, Status: types.Present},
		&types.Interval{Microseconds: -1000001, Status: types.Present},
		&types.Interval{Microseconds: -123202800000000, Status: types.Present},
		&types.Interval{Days: -1, Status: types.Present},
		&types.Interval{Months: -1, Status: types.Present},
		&types.Interval{Months: -12, Status: types.Present},
		&types.Interval{Months: -13, Days: -15, Microseconds: -1000001, Status: types.Present},
		&types.Interval{Status: types.Null},
	})
}

func TestIntervalNormalize(t *testing.T) {
	testutil.TestSuccessfulNormalize(t, []testutil.NormalizeTest{
		{
			SQL:   "select '1 second'::interval",
			Value: &types.Interval{Microseconds: 1000000, Status: types.Present},
		},
		{
			SQL:   "select '1.000001 second'::interval",
			Value: &types.Interval{Microseconds: 1000001, Status: types.Present},
		},
		{
			SQL:   "select '34223 hours'::interval",
			Value: &types.Interval{Microseconds: 123202800000000, Status: types.Present},
		},
		{
			SQL:   "select '1 day'::interval",
			Value: &types.Interval{Days: 1, Status: types.Present},
		},
		{
			SQL:   "select '1 month'::interval",
			Value: &types.Interval{Months: 1, Status: types.Present},
		},
		{
			SQL:   "select '1 year'::interval",
			Value: &types.Interval{Months: 12, Status: types.Present},
		},
		{
			SQL:   "select '-13 mon'::interval",
			Value: &types.Interval{Months: -13, Status: types.Present},
		},
	})
}
