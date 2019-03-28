package types_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/elliotcourant/noahdb/pkg/types"
	"github.com/elliotcourant/noahdb/pkg/types/testutil"
)

func TestTimestampArrayTranscode(t *testing.T) {
	testutil.TestSuccessfulTranscodeEqFunc(t, "timestamp[]", []interface{}{
		&types.TimestampArray{
			Elements:   nil,
			Dimensions: nil,
			Status:     types.Present,
		},
		&types.TimestampArray{
			Elements: []types.Timestamp{
				{Time: time.Date(2015, 2, 1, 0, 0, 0, 0, time.UTC), Status: types.Present},
				{Status: types.Null},
			},
			Dimensions: []types.ArrayDimension{{Length: 2, LowerBound: 1}},
			Status:     types.Present,
		},
		&types.TimestampArray{Status: types.Null},
		&types.TimestampArray{
			Elements: []types.Timestamp{
				{Time: time.Date(2015, 2, 1, 0, 0, 0, 0, time.UTC), Status: types.Present},
				{Time: time.Date(2016, 2, 1, 0, 0, 0, 0, time.UTC), Status: types.Present},
				{Time: time.Date(2017, 2, 1, 0, 0, 0, 0, time.UTC), Status: types.Present},
				{Time: time.Date(2012, 1, 1, 0, 0, 0, 0, time.UTC), Status: types.Present},
				{Status: types.Null},
				{Time: time.Date(2015, 2, 1, 0, 0, 0, 0, time.UTC), Status: types.Present},
			},
			Dimensions: []types.ArrayDimension{{Length: 3, LowerBound: 1}, {Length: 2, LowerBound: 1}},
			Status:     types.Present,
		},
		&types.TimestampArray{
			Elements: []types.Timestamp{
				{Time: time.Date(2015, 2, 1, 0, 0, 0, 0, time.UTC), Status: types.Present},
				{Time: time.Date(2015, 2, 2, 0, 0, 0, 0, time.UTC), Status: types.Present},
				{Time: time.Date(2015, 2, 3, 0, 0, 0, 0, time.UTC), Status: types.Present},
				{Time: time.Date(2015, 2, 4, 0, 0, 0, 0, time.UTC), Status: types.Present},
			},
			Dimensions: []types.ArrayDimension{
				{Length: 2, LowerBound: 4},
				{Length: 2, LowerBound: 2},
			},
			Status: types.Present,
		},
	}, func(a, b interface{}) bool {
		ata := a.(types.TimestampArray)
		bta := b.(types.TimestampArray)

		if len(ata.Elements) != len(bta.Elements) || ata.Status != bta.Status {
			return false
		}

		for i := range ata.Elements {
			ae, be := ata.Elements[i], bta.Elements[i]
			if !(ae.Time.Equal(be.Time) && ae.Status == be.Status && ae.InfinityModifier == be.InfinityModifier) {
				return false
			}
		}

		return true
	})
}

func TestTimestampArraySet(t *testing.T) {
	successfulTests := []struct {
		source interface{}
		result types.TimestampArray
	}{
		{
			source: []time.Time{time.Date(2015, 2, 1, 0, 0, 0, 0, time.UTC)},
			result: types.TimestampArray{
				Elements:   []types.Timestamp{{Time: time.Date(2015, 2, 1, 0, 0, 0, 0, time.UTC), Status: types.Present}},
				Dimensions: []types.ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     types.Present},
		},
		{
			source: ([]time.Time)(nil),
			result: types.TimestampArray{Status: types.Null},
		},
	}

	for i, tt := range successfulTests {
		var r types.TimestampArray
		err := r.Set(tt.source)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if !reflect.DeepEqual(r, tt.result) {
			t.Errorf("%d: expected %v to convert to %v, but it was %v", i, tt.source, tt.result, r)
		}
	}
}

func TestTimestampArrayAssignTo(t *testing.T) {
	var timeSlice []time.Time

	simpleTests := []struct {
		src      types.TimestampArray
		dst      interface{}
		expected interface{}
	}{
		{
			src: types.TimestampArray{
				Elements:   []types.Timestamp{{Time: time.Date(2015, 2, 1, 0, 0, 0, 0, time.UTC), Status: types.Present}},
				Dimensions: []types.ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     types.Present,
			},
			dst:      &timeSlice,
			expected: []time.Time{time.Date(2015, 2, 1, 0, 0, 0, 0, time.UTC)},
		},
		{
			src:      types.TimestampArray{Status: types.Null},
			dst:      &timeSlice,
			expected: ([]time.Time)(nil),
		},
	}

	for i, tt := range simpleTests {
		err := tt.src.AssignTo(tt.dst)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if dst := reflect.ValueOf(tt.dst).Elem().Interface(); !reflect.DeepEqual(dst, tt.expected) {
			t.Errorf("%d: expected %v to assign %v, but result was %v", i, tt.src, tt.expected, dst)
		}
	}

	errorTests := []struct {
		src types.TimestampArray
		dst interface{}
	}{
		{
			src: types.TimestampArray{
				Elements:   []types.Timestamp{{Status: types.Null}},
				Dimensions: []types.ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     types.Present,
			},
			dst: &timeSlice,
		},
	}

	for i, tt := range errorTests {
		err := tt.src.AssignTo(tt.dst)
		if err == nil {
			t.Errorf("%d: expected error but none was returned (%v -> %v)", i, tt.src, tt.dst)
		}
	}

}
