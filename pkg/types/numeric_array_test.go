package types_test

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/elliotcourant/noahdb/pkg/types"
	"github.com/elliotcourant/noahdb/pkg/types/testutil"
)

func TestNumericArrayTranscode(t *testing.T) {
	testutil.TestSuccessfulTranscode(t, "numeric[]", []interface{}{
		&types.NumericArray{
			Elements:   nil,
			Dimensions: nil,
			Status:     types.Present,
		},
		&types.NumericArray{
			Elements: []types.Numeric{
				{Int: big.NewInt(1), Status: types.Present},
				{Status: types.Null},
			},
			Dimensions: []types.ArrayDimension{{Length: 2, LowerBound: 1}},
			Status:     types.Present,
		},
		&types.NumericArray{Status: types.Null},
		&types.NumericArray{
			Elements: []types.Numeric{
				{Int: big.NewInt(1), Status: types.Present},
				{Int: big.NewInt(2), Status: types.Present},
				{Int: big.NewInt(3), Status: types.Present},
				{Int: big.NewInt(4), Status: types.Present},
				{Status: types.Null},
				{Int: big.NewInt(6), Status: types.Present},
			},
			Dimensions: []types.ArrayDimension{{Length: 3, LowerBound: 1}, {Length: 2, LowerBound: 1}},
			Status:     types.Present,
		},
		&types.NumericArray{
			Elements: []types.Numeric{
				{Int: big.NewInt(1), Status: types.Present},
				{Int: big.NewInt(2), Status: types.Present},
				{Int: big.NewInt(3), Status: types.Present},
				{Int: big.NewInt(4), Status: types.Present},
			},
			Dimensions: []types.ArrayDimension{
				{Length: 2, LowerBound: 4},
				{Length: 2, LowerBound: 2},
			},
			Status: types.Present,
		},
	})
}

func TestNumericArraySet(t *testing.T) {
	successfulTests := []struct {
		source interface{}
		result types.NumericArray
	}{
		{
			source: []float32{1},
			result: types.NumericArray{
				Elements:   []types.Numeric{{Int: big.NewInt(1), Status: types.Present}},
				Dimensions: []types.ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     types.Present},
		},
		{
			source: []float64{1},
			result: types.NumericArray{
				Elements:   []types.Numeric{{Int: big.NewInt(1), Status: types.Present}},
				Dimensions: []types.ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     types.Present},
		},
		{
			source: ([]float32)(nil),
			result: types.NumericArray{Status: types.Null},
		},
	}

	for i, tt := range successfulTests {
		var r types.NumericArray
		err := r.Set(tt.source)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if !reflect.DeepEqual(r, tt.result) {
			t.Errorf("%d: expected %v to convert to %v, but it was %v", i, tt.source, tt.result, r)
		}
	}
}

func TestNumericArrayAssignTo(t *testing.T) {
	var float32Slice []float32
	var float64Slice []float64

	simpleTests := []struct {
		src      types.NumericArray
		dst      interface{}
		expected interface{}
	}{
		{
			src: types.NumericArray{
				Elements:   []types.Numeric{{Int: big.NewInt(1), Status: types.Present}},
				Dimensions: []types.ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     types.Present,
			},
			dst:      &float32Slice,
			expected: []float32{1},
		},
		{
			src: types.NumericArray{
				Elements:   []types.Numeric{{Int: big.NewInt(1), Status: types.Present}},
				Dimensions: []types.ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     types.Present,
			},
			dst:      &float64Slice,
			expected: []float64{1},
		},
		{
			src:      types.NumericArray{Status: types.Null},
			dst:      &float32Slice,
			expected: ([]float32)(nil),
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
		src types.NumericArray
		dst interface{}
	}{
		{
			src: types.NumericArray{
				Elements:   []types.Numeric{{Status: types.Null}},
				Dimensions: []types.ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     types.Present,
			},
			dst: &float32Slice,
		},
	}

	for i, tt := range errorTests {
		err := tt.src.AssignTo(tt.dst)
		if err == nil {
			t.Errorf("%d: expected error but none was returned (%v -> %v)", i, tt.src, tt.dst)
		}
	}

}
