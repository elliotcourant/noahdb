package types_test

import (
	"math"
	"reflect"
	"testing"

	"github.com/elliotcourant/noahdb/pkg/types"
	"github.com/elliotcourant/noahdb/pkg/types/testutil"
)

func TestInt2Transcode(t *testing.T) {
	testutil.TestSuccessfulTranscode(t, "int2", []interface{}{
		&types.Int2{Int: math.MinInt16, Status: types.Present},
		&types.Int2{Int: -1, Status: types.Present},
		&types.Int2{Int: 0, Status: types.Present},
		&types.Int2{Int: 1, Status: types.Present},
		&types.Int2{Int: math.MaxInt16, Status: types.Present},
		&types.Int2{Int: 0, Status: types.Null},
	})
}

func TestInt2Set(t *testing.T) {
	successfulTests := []struct {
		source interface{}
		result types.Int2
	}{
		{source: int8(1), result: types.Int2{Int: 1, Status: types.Present}},
		{source: int16(1), result: types.Int2{Int: 1, Status: types.Present}},
		{source: int32(1), result: types.Int2{Int: 1, Status: types.Present}},
		{source: int64(1), result: types.Int2{Int: 1, Status: types.Present}},
		{source: int8(-1), result: types.Int2{Int: -1, Status: types.Present}},
		{source: int16(-1), result: types.Int2{Int: -1, Status: types.Present}},
		{source: int32(-1), result: types.Int2{Int: -1, Status: types.Present}},
		{source: int64(-1), result: types.Int2{Int: -1, Status: types.Present}},
		{source: uint8(1), result: types.Int2{Int: 1, Status: types.Present}},
		{source: uint16(1), result: types.Int2{Int: 1, Status: types.Present}},
		{source: uint32(1), result: types.Int2{Int: 1, Status: types.Present}},
		{source: uint64(1), result: types.Int2{Int: 1, Status: types.Present}},
		{source: "1", result: types.Int2{Int: 1, Status: types.Present}},
		{source: _int8(1), result: types.Int2{Int: 1, Status: types.Present}},
	}

	for i, tt := range successfulTests {
		var r types.Int2
		err := r.Set(tt.source)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if r != tt.result {
			t.Errorf("%d: expected %v to convert to %v, but it was %v", i, tt.source, tt.result, r)
		}
	}
}

func TestInt2AssignTo(t *testing.T) {
	var i8 int8
	var i16 int16
	var i32 int32
	var i64 int64
	var i int
	var ui8 uint8
	var ui16 uint16
	var ui32 uint32
	var ui64 uint64
	var ui uint
	var pi8 *int8
	var _i8 _int8
	var _pi8 *_int8

	simpleTests := []struct {
		src      types.Int2
		dst      interface{}
		expected interface{}
	}{
		{src: types.Int2{Int: 42, Status: types.Present}, dst: &i8, expected: int8(42)},
		{src: types.Int2{Int: 42, Status: types.Present}, dst: &i16, expected: int16(42)},
		{src: types.Int2{Int: 42, Status: types.Present}, dst: &i32, expected: int32(42)},
		{src: types.Int2{Int: 42, Status: types.Present}, dst: &i64, expected: int64(42)},
		{src: types.Int2{Int: 42, Status: types.Present}, dst: &i, expected: int(42)},
		{src: types.Int2{Int: 42, Status: types.Present}, dst: &ui8, expected: uint8(42)},
		{src: types.Int2{Int: 42, Status: types.Present}, dst: &ui16, expected: uint16(42)},
		{src: types.Int2{Int: 42, Status: types.Present}, dst: &ui32, expected: uint32(42)},
		{src: types.Int2{Int: 42, Status: types.Present}, dst: &ui64, expected: uint64(42)},
		{src: types.Int2{Int: 42, Status: types.Present}, dst: &ui, expected: uint(42)},
		{src: types.Int2{Int: 42, Status: types.Present}, dst: &_i8, expected: _int8(42)},
		{src: types.Int2{Int: 0, Status: types.Null}, dst: &pi8, expected: (*int8)(nil)},
		{src: types.Int2{Int: 0, Status: types.Null}, dst: &_pi8, expected: (*_int8)(nil)},
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
		src      types.Int2
		dst      interface{}
		expected interface{}
	}{
		{src: types.Int2{Int: 42, Status: types.Present}, dst: &pi8, expected: int8(42)},
		{src: types.Int2{Int: 42, Status: types.Present}, dst: &_pi8, expected: _int8(42)},
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
		src types.Int2
		dst interface{}
	}{
		{src: types.Int2{Int: 150, Status: types.Present}, dst: &i8},
		{src: types.Int2{Int: -1, Status: types.Present}, dst: &ui8},
		{src: types.Int2{Int: -1, Status: types.Present}, dst: &ui16},
		{src: types.Int2{Int: -1, Status: types.Present}, dst: &ui32},
		{src: types.Int2{Int: -1, Status: types.Present}, dst: &ui64},
		{src: types.Int2{Int: -1, Status: types.Present}, dst: &ui},
		{src: types.Int2{Int: 0, Status: types.Null}, dst: &i16},
	}

	for i, tt := range errorTests {
		err := tt.src.AssignTo(tt.dst)
		if err == nil {
			t.Errorf("%d: expected error but none was returned (%v -> %v)", i, tt.src, tt.dst)
		}
	}
}
