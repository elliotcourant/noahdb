package types_test

import (
	"reflect"
	"testing"

	"github.com/elliotcourant/noahdb/pkg/types"
	"github.com/elliotcourant/noahdb/pkg/types/testutil"
)

func TestInt2ArrayTranscode(t *testing.T) {
	testutil.TestSuccessfulTranscode(t, "int2[]", []interface{}{
		&types.Int2Array{
			Elements:   nil,
			Dimensions: nil,
			Status:     types.Present,
		},
		&types.Int2Array{
			Elements: []types.Int2{
				{Int: 1, Status: types.Present},
				{Status: types.Null},
			},
			Dimensions: []types.ArrayDimension{{Length: 2, LowerBound: 1}},
			Status:     types.Present,
		},
		&types.Int2Array{Status: types.Null},
		&types.Int2Array{
			Elements: []types.Int2{
				{Int: 1, Status: types.Present},
				{Int: 2, Status: types.Present},
				{Int: 3, Status: types.Present},
				{Int: 4, Status: types.Present},
				{Status: types.Null},
				{Int: 6, Status: types.Present},
			},
			Dimensions: []types.ArrayDimension{{Length: 3, LowerBound: 1}, {Length: 2, LowerBound: 1}},
			Status:     types.Present,
		},
		&types.Int2Array{
			Elements: []types.Int2{
				{Int: 1, Status: types.Present},
				{Int: 2, Status: types.Present},
				{Int: 3, Status: types.Present},
				{Int: 4, Status: types.Present},
			},
			Dimensions: []types.ArrayDimension{
				{Length: 2, LowerBound: 4},
				{Length: 2, LowerBound: 2},
			},
			Status: types.Present,
		},
	})
}

func TestInt2ArraySet(t *testing.T) {
	successfulTests := []struct {
		source interface{}
		result types.Int2Array
	}{
		{
			source: []int16{1},
			result: types.Int2Array{
				Elements:   []types.Int2{{Int: 1, Status: types.Present}},
				Dimensions: []types.ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     types.Present},
		},
		{
			source: []uint16{1},
			result: types.Int2Array{
				Elements:   []types.Int2{{Int: 1, Status: types.Present}},
				Dimensions: []types.ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     types.Present},
		},
		{
			source: ([]int16)(nil),
			result: types.Int2Array{Status: types.Null},
		},
	}

	for i, tt := range successfulTests {
		var r types.Int2Array
		err := r.Set(tt.source)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if !reflect.DeepEqual(r, tt.result) {
			t.Errorf("%d: expected %v to convert to %v, but it was %v", i, tt.source, tt.result, r)
		}
	}
}

func TestInt2ArrayAssignTo(t *testing.T) {
	var int16Slice []int16
	var uint16Slice []uint16
	var namedInt16Slice _int16Slice

	simpleTests := []struct {
		src      types.Int2Array
		dst      interface{}
		expected interface{}
	}{
		{
			src: types.Int2Array{
				Elements:   []types.Int2{{Int: 1, Status: types.Present}},
				Dimensions: []types.ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     types.Present,
			},
			dst:      &int16Slice,
			expected: []int16{1},
		},
		{
			src: types.Int2Array{
				Elements:   []types.Int2{{Int: 1, Status: types.Present}},
				Dimensions: []types.ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     types.Present,
			},
			dst:      &uint16Slice,
			expected: []uint16{1},
		},
		{
			src: types.Int2Array{
				Elements:   []types.Int2{{Int: 1, Status: types.Present}},
				Dimensions: []types.ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     types.Present,
			},
			dst:      &namedInt16Slice,
			expected: _int16Slice{1},
		},
		{
			src:      types.Int2Array{Status: types.Null},
			dst:      &int16Slice,
			expected: ([]int16)(nil),
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
		src types.Int2Array
		dst interface{}
	}{
		{
			src: types.Int2Array{
				Elements:   []types.Int2{{Status: types.Null}},
				Dimensions: []types.ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     types.Present,
			},
			dst: &int16Slice,
		},
		{
			src: types.Int2Array{
				Elements:   []types.Int2{{Int: -1, Status: types.Present}},
				Dimensions: []types.ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     types.Present,
			},
			dst: &uint16Slice,
		},
	}

	for i, tt := range errorTests {
		err := tt.src.AssignTo(tt.dst)
		if err == nil {
			t.Errorf("%d: expected error but none was returned (%v -> %v)", i, tt.src, tt.dst)
		}
	}

}
