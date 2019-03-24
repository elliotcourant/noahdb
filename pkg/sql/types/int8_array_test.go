/*
 * Copyright (c) 2019 Ready Stock
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing
 * permissions and limitations under the License.
 */

package types_test

import (
	"reflect"
	"testing"

	"github.com/elliotcourant/noahdb/pkg/sql/types"
	"github.com/elliotcourant/noahdb/pkg/sql/types/testutil"
)

func TestInt8ArrayTranscode(t *testing.T) {
	testutil.TestSuccessfulTranscode(t, "int8[]", []interface{}{
		&types.Int8Array{
			Elements:   nil,
			Dimensions: nil,
			Status:     types.Present,
		},
		&types.Int8Array{
			Elements: []types.Int8{
				{Int: 1, Status: types.Present},
				{Status: types.Null},
			},
			Dimensions: []types.ArrayDimension{{Length: 2, LowerBound: 1}},
			Status:     types.Present,
		},
		&types.Int8Array{Status: types.Null},
		&types.Int8Array{
			Elements: []types.Int8{
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
		&types.Int8Array{
			Elements: []types.Int8{
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

func TestInt8ArraySet(t *testing.T) {
	successfulTests := []struct {
		source interface{}
		result types.Int8Array
	}{
		{
			source: []int64{1},
			result: types.Int8Array{
				Elements:   []types.Int8{{Int: 1, Status: types.Present}},
				Dimensions: []types.ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     types.Present},
		},
		{
			source: []uint64{1},
			result: types.Int8Array{
				Elements:   []types.Int8{{Int: 1, Status: types.Present}},
				Dimensions: []types.ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     types.Present},
		},
		{
			source: ([]int64)(nil),
			result: types.Int8Array{Status: types.Null},
		},
	}

	for i, tt := range successfulTests {
		var r types.Int8Array
		err := r.Set(tt.source)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if !reflect.DeepEqual(r, tt.result) {
			t.Errorf("%d: expected %v to convert to %v, but it was %v", i, tt.source, tt.result, r)
		}
	}
}

func TestInt8ArrayAssignTo(t *testing.T) {
	var int64Slice []int64
	var uint64Slice []uint64
	var namedInt64Slice _int64Slice

	simpleTests := []struct {
		src      types.Int8Array
		dst      interface{}
		expected interface{}
	}{
		{
			src: types.Int8Array{
				Elements:   []types.Int8{{Int: 1, Status: types.Present}},
				Dimensions: []types.ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     types.Present,
			},
			dst:      &int64Slice,
			expected: []int64{1},
		},
		{
			src: types.Int8Array{
				Elements:   []types.Int8{{Int: 1, Status: types.Present}},
				Dimensions: []types.ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     types.Present,
			},
			dst:      &uint64Slice,
			expected: []uint64{1},
		},
		{
			src: types.Int8Array{
				Elements:   []types.Int8{{Int: 1, Status: types.Present}},
				Dimensions: []types.ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     types.Present,
			},
			dst:      &namedInt64Slice,
			expected: _int64Slice{1},
		},
		{
			src:      types.Int8Array{Status: types.Null},
			dst:      &int64Slice,
			expected: ([]int64)(nil),
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
		src types.Int8Array
		dst interface{}
	}{
		{
			src: types.Int8Array{
				Elements:   []types.Int8{{Status: types.Null}},
				Dimensions: []types.ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     types.Present,
			},
			dst: &int64Slice,
		},
		{
			src: types.Int8Array{
				Elements:   []types.Int8{{Int: -1, Status: types.Present}},
				Dimensions: []types.ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     types.Present,
			},
			dst: &uint64Slice,
		},
	}

	for i, tt := range errorTests {
		err := tt.src.AssignTo(tt.dst)
		if err == nil {
			t.Errorf("%d: expected error but none was returned (%v -> %v)", i, tt.src, tt.dst)
		}
	}

}
