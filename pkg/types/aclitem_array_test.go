package types_test

import (
	"reflect"
	"testing"

	"github.com/elliotcourant/noahdb/pkg/types"
	"github.com/elliotcourant/noahdb/pkg/types/testutil"
)

func TestACLItemArrayTranscode(t *testing.T) {
	testutil.TestSuccessfulTranscode(t, "aclitem[]", []interface{}{
		&types.ACLItemArray{
			Elements:   nil,
			Dimensions: nil,
			Status:     types.Present,
		},
		&types.ACLItemArray{
			Elements: []types.ACLItem{
				{String: "=r/postgres", Status: types.Present},
				{Status: types.Null},
			},
			Dimensions: []types.ArrayDimension{{Length: 2, LowerBound: 1}},
			Status:     types.Present,
		},
		&types.ACLItemArray{Status: types.Null},
		&types.ACLItemArray{
			Elements: []types.ACLItem{
				{String: "=r/postgres", Status: types.Present},
				{String: "postgres=arwdDxt/postgres", Status: types.Present},
				{String: `postgres=arwdDxt/" tricky, ' } "" \ test user "`, Status: types.Present},
				{String: "=r/postgres", Status: types.Present},
				{Status: types.Null},
				{String: "=r/postgres", Status: types.Present},
			},
			Dimensions: []types.ArrayDimension{{Length: 3, LowerBound: 1}, {Length: 2, LowerBound: 1}},
			Status:     types.Present,
		},
		&types.ACLItemArray{
			Elements: []types.ACLItem{
				{String: "=r/postgres", Status: types.Present},
				{String: "postgres=arwdDxt/postgres", Status: types.Present},
				{String: "=r/postgres", Status: types.Present},
				{String: "postgres=arwdDxt/postgres", Status: types.Present},
			},
			Dimensions: []types.ArrayDimension{
				{Length: 2, LowerBound: 4},
				{Length: 2, LowerBound: 2},
			},
			Status: types.Present,
		},
	})
}

func TestACLItemArraySet(t *testing.T) {
	successfulTests := []struct {
		source interface{}
		result types.ACLItemArray
	}{
		{
			source: []string{"=r/postgres"},
			result: types.ACLItemArray{
				Elements:   []types.ACLItem{{String: "=r/postgres", Status: types.Present}},
				Dimensions: []types.ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     types.Present},
		},
		{
			source: ([]string)(nil),
			result: types.ACLItemArray{Status: types.Null},
		},
	}

	for i, tt := range successfulTests {
		var r types.ACLItemArray
		err := r.Set(tt.source)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if !reflect.DeepEqual(r, tt.result) {
			t.Errorf("%d: expected %v to convert to %v, but it was %v", i, tt.source, tt.result, r)
		}
	}
}

func TestACLItemArrayAssignTo(t *testing.T) {
	var stringSlice []string
	type _stringSlice []string
	var namedStringSlice _stringSlice

	simpleTests := []struct {
		src      types.ACLItemArray
		dst      interface{}
		expected interface{}
	}{
		{
			src: types.ACLItemArray{
				Elements:   []types.ACLItem{{String: "=r/postgres", Status: types.Present}},
				Dimensions: []types.ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     types.Present,
			},
			dst:      &stringSlice,
			expected: []string{"=r/postgres"},
		},
		{
			src: types.ACLItemArray{
				Elements:   []types.ACLItem{{String: "=r/postgres", Status: types.Present}},
				Dimensions: []types.ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     types.Present,
			},
			dst:      &namedStringSlice,
			expected: _stringSlice{"=r/postgres"},
		},
		{
			src:      types.ACLItemArray{Status: types.Null},
			dst:      &stringSlice,
			expected: ([]string)(nil),
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
		src types.ACLItemArray
		dst interface{}
	}{
		{
			src: types.ACLItemArray{
				Elements:   []types.ACLItem{{Status: types.Null}},
				Dimensions: []types.ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     types.Present,
			},
			dst: &stringSlice,
		},
	}

	for i, tt := range errorTests {
		err := tt.src.AssignTo(tt.dst)
		if err == nil {
			t.Errorf("%d: expected error but none was returned (%v -> %v)", i, tt.src, tt.dst)
		}
	}
}
