package types_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/elliotcourant/noahdb/pkg/types"
	"github.com/elliotcourant/noahdb/pkg/types/testutil"
)

func TestTimestampTranscode(t *testing.T) {
	testutil.TestSuccessfulTranscodeEqFunc(t, "timestamp", []interface{}{
		&types.Timestamp{Time: time.Date(1800, 1, 1, 0, 0, 0, 0, time.UTC), Status: types.Present},
		&types.Timestamp{Time: time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC), Status: types.Present},
		&types.Timestamp{Time: time.Date(1905, 1, 1, 0, 0, 0, 0, time.UTC), Status: types.Present},
		&types.Timestamp{Time: time.Date(1940, 1, 1, 0, 0, 0, 0, time.UTC), Status: types.Present},
		&types.Timestamp{Time: time.Date(1960, 1, 1, 0, 0, 0, 0, time.UTC), Status: types.Present},
		&types.Timestamp{Time: time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC), Status: types.Present},
		&types.Timestamp{Time: time.Date(1999, 12, 31, 0, 0, 0, 0, time.UTC), Status: types.Present},
		&types.Timestamp{Time: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), Status: types.Present},
		&types.Timestamp{Time: time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC), Status: types.Present},
		&types.Timestamp{Time: time.Date(2200, 1, 1, 0, 0, 0, 0, time.UTC), Status: types.Present},
		&types.Timestamp{Status: types.Null},
		&types.Timestamp{Status: types.Present, InfinityModifier: types.Infinity},
		&types.Timestamp{Status: types.Present, InfinityModifier: -types.Infinity},
	}, func(a, b interface{}) bool {
		at := a.(types.Timestamp)
		bt := b.(types.Timestamp)

		return at.Time.Equal(bt.Time) && at.Status == bt.Status && at.InfinityModifier == bt.InfinityModifier
	})
}

func TestTimestampSet(t *testing.T) {
	type _time time.Time

	successfulTests := []struct {
		source interface{}
		result types.Timestamp
	}{
		{source: time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC), result: types.Timestamp{Time: time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC), Status: types.Present}},
		{source: time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC), result: types.Timestamp{Time: time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC), Status: types.Present}},
		{source: time.Date(1999, 12, 31, 12, 59, 59, 0, time.UTC), result: types.Timestamp{Time: time.Date(1999, 12, 31, 12, 59, 59, 0, time.UTC), Status: types.Present}},
		{source: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), result: types.Timestamp{Time: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), Status: types.Present}},
		{source: time.Date(2000, 1, 1, 0, 0, 1, 0, time.UTC), result: types.Timestamp{Time: time.Date(2000, 1, 1, 0, 0, 1, 0, time.UTC), Status: types.Present}},
		{source: time.Date(2200, 1, 1, 0, 0, 0, 0, time.UTC), result: types.Timestamp{Time: time.Date(2200, 1, 1, 0, 0, 0, 0, time.UTC), Status: types.Present}},
		{source: time.Date(2015, 1, 1, 0, 0, 0, 0, time.Local), result: types.Timestamp{Time: time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC), Status: types.Present}},
		{source: _time(time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)), result: types.Timestamp{Time: time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC), Status: types.Present}},
	}

	for i, tt := range successfulTests {
		var r types.Timestamp
		err := r.Set(tt.source)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if r != tt.result {
			t.Errorf("%d: expected %v to convert to %v, but it was %v", i, tt.source, tt.result, r)
		}
	}
}

func TestTimestampAssignTo(t *testing.T) {
	var tim time.Time
	var ptim *time.Time

	simpleTests := []struct {
		src      types.Timestamp
		dst      interface{}
		expected interface{}
	}{
		{src: types.Timestamp{Time: time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC), Status: types.Present}, dst: &tim, expected: time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)},
		{src: types.Timestamp{Time: time.Time{}, Status: types.Null}, dst: &ptim, expected: (*time.Time)(nil)},
	}

	for i, tt := range simpleTests {
		err := tt.src.AssignTo(tt.dst)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if dst := reflect.ValueOf(tt.dst).Elem().Interface(); dst != tt.expected {
			t.Errorf("%d: expected %v to assign %v, but result was %v", i, tt.src, tt.expected, dst)
		}
	}

	pointerAllocTests := []struct {
		src      types.Timestamp
		dst      interface{}
		expected interface{}
	}{
		{src: types.Timestamp{Time: time.Date(2015, 1, 1, 0, 0, 0, 0, time.Local), Status: types.Present}, dst: &ptim, expected: time.Date(2015, 1, 1, 0, 0, 0, 0, time.Local)},
	}

	for i, tt := range pointerAllocTests {
		err := tt.src.AssignTo(tt.dst)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if dst := reflect.ValueOf(tt.dst).Elem().Elem().Interface(); dst != tt.expected {
			t.Errorf("%d: expected %v to assign %v, but result was %v", i, tt.src, tt.expected, dst)
		}
	}

	errorTests := []struct {
		src types.Timestamp
		dst interface{}
	}{
		{src: types.Timestamp{Time: time.Date(2015, 1, 1, 0, 0, 0, 0, time.Local), InfinityModifier: types.Infinity, Status: types.Present}, dst: &tim},
		{src: types.Timestamp{Time: time.Date(2015, 1, 1, 0, 0, 0, 0, time.Local), InfinityModifier: types.NegativeInfinity, Status: types.Present}, dst: &tim},
		{src: types.Timestamp{Time: time.Date(2015, 1, 1, 0, 0, 0, 0, time.Local), Status: types.Null}, dst: &tim},
	}

	for i, tt := range errorTests {
		err := tt.src.AssignTo(tt.dst)
		if err == nil {
			t.Errorf("%d: expected error but none was returned (%v -> %v)", i, tt.src, tt.dst)
		}
	}
}
