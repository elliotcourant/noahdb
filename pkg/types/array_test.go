package types_test

import (
	"reflect"
	"testing"

	"github.com/elliotcourant/noahdb/pkg/types"
)

func TestParseUntypedTextArray(t *testing.T) {
	tests := []struct {
		source string
		result types.UntypedTextArray
	}{
		{
			source: "{}",
			result: types.UntypedTextArray{
				Elements:   nil,
				Dimensions: nil,
			},
		},
		{
			source: "{1}",
			result: types.UntypedTextArray{
				Elements:   []string{"1"},
				Dimensions: []types.ArrayDimension{{Length: 1, LowerBound: 1}},
			},
		},
		{
			source: "{a,b}",
			result: types.UntypedTextArray{
				Elements:   []string{"a", "b"},
				Dimensions: []types.ArrayDimension{{Length: 2, LowerBound: 1}},
			},
		},
		{
			source: `{"NULL"}`,
			result: types.UntypedTextArray{
				Elements:   []string{"NULL"},
				Dimensions: []types.ArrayDimension{{Length: 1, LowerBound: 1}},
			},
		},
		{
			source: `{""}`,
			result: types.UntypedTextArray{
				Elements:   []string{""},
				Dimensions: []types.ArrayDimension{{Length: 1, LowerBound: 1}},
			},
		},
		{
			source: `{"He said, \"Hello.\""}`,
			result: types.UntypedTextArray{
				Elements:   []string{`He said, "Hello."`},
				Dimensions: []types.ArrayDimension{{Length: 1, LowerBound: 1}},
			},
		},
		{
			source: "{{a,b},{c,d},{e,f}}",
			result: types.UntypedTextArray{
				Elements:   []string{"a", "b", "c", "d", "e", "f"},
				Dimensions: []types.ArrayDimension{{Length: 3, LowerBound: 1}, {Length: 2, LowerBound: 1}},
			},
		},
		{
			source: "{{{a,b},{c,d},{e,f}},{{a,b},{c,d},{e,f}}}",
			result: types.UntypedTextArray{
				Elements: []string{"a", "b", "c", "d", "e", "f", "a", "b", "c", "d", "e", "f"},
				Dimensions: []types.ArrayDimension{
					{Length: 2, LowerBound: 1},
					{Length: 3, LowerBound: 1},
					{Length: 2, LowerBound: 1},
				},
			},
		},
		{
			source: "[4:4]={1}",
			result: types.UntypedTextArray{
				Elements:   []string{"1"},
				Dimensions: []types.ArrayDimension{{Length: 1, LowerBound: 4}},
			},
		},
		{
			source: "[4:5][2:3]={{a,b},{c,d}}",
			result: types.UntypedTextArray{
				Elements: []string{"a", "b", "c", "d"},
				Dimensions: []types.ArrayDimension{
					{Length: 2, LowerBound: 4},
					{Length: 2, LowerBound: 2},
				},
			},
		},
	}

	for i, tt := range tests {
		r, err := types.ParseUntypedTextArray(tt.source)
		if err != nil {
			t.Errorf("%d: %v", i, err)
			continue
		}

		if !reflect.DeepEqual(*r, tt.result) {
			t.Errorf("%d: expected %+v to be parsed to %+v, but it was %+v", i, tt.source, tt.result, *r)
		}
	}
}
