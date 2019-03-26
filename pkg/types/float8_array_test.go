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

	"github.com/elliotcourant/noahdb/pkg/types"
	"github.com/elliotcourant/noahdb/pkg/types/testutil"
)

func TestFloat8ArrayTranscode(t *testing.T) {
	testutil.TestSuccessfulTranscode(t, "float8[]", []interface{}{
		&types.Float8Array{
			Elements:   nil,
			Dimensions: nil,
			Status:     types.Present,
		},
		&types.Float8Array{
			Elements: []types.Float8{
				{Float: 1, Status: types.Present},
				{Status: types.Null},
			},
			Dimensions: []types.ArrayDimension{{Length: 2, LowerBound: 1}},
			Status:     types.Present,
		},
		&types.Float8Array{Status: types.Null},
		&types.Float8Array{
			Elements: []types.Float8{
				{Float: 1, Status: types.Present},
				{Float: 2, Status: types.Present},
				{Float: 3, Status: types.Present},
				{Float: 4, Status: types.Present},
				{Status: types.Null},
				{Float: 6, Status: types.Present},
			},
			Dimensions: []types.ArrayDimension{{Length: 3, LowerBound: 1}, {Length: 2, LowerBound: 1}},
			Status:     types.Present,
		},
		&types.Float8Array{
			Elements: []types.Float8{
				{Float: 1, Status: types.Present},
				{Float: 2, Status: types.Present},
				{Float: 3, Status: types.Present},
				{Float: 4, Status: types.Present},
			},
			Dimensions: []types.ArrayDimension{
				{Length: 2, LowerBound: 4},
				{Length: 2, LowerBound: 2},
			},
			Status: types.Present,
		},
	})
}

func TestFloat8ArraySet(t *testing.T) {
	successfulTests := []struct {
		source interface{}
		result types.Float8Array
	}{
		{
			source: []float64{1},
			result: types.Float8Array{
				Elements:   []types.Float8{{Float: 1, Status: types.Present}},
				Dimensions: []types.ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     types.Present},
		},
		{
			source: ([]float64)(nil),
			result: types.Float8Array{Status: types.Null},
		},
	}

	for i, tt := range successfulTests {
		var r types.Float8Array
		err := r.Set(tt.source)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if !reflect.DeepEqual(r, tt.result) {
			t.Errorf("%d: expected %v to convert to %v, but it was %v", i, tt.source, tt.result, r)
		}
	}
}

func TestFloat8ArrayAssignTo(t *testing.T) {
	var float64Slice []float64
	var namedFloat64Slice _float64Slice

	simpleTests := []struct {
		src      types.Float8Array
		dst      interface{}
		expected interface{}
	}{
		{
			src: types.Float8Array{
				Elements:   []types.Float8{{Float: 1.23, Status: types.Present}},
				Dimensions: []types.ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     types.Present,
			},
			dst:      &float64Slice,
			expected: []float64{1.23},
		},
		{
			src: types.Float8Array{
				Elements:   []types.Float8{{Float: 1.23, Status: types.Present}},
				Dimensions: []types.ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     types.Present,
			},
			dst:      &namedFloat64Slice,
			expected: _float64Slice{1.23},
		},
		{
			src:      types.Float8Array{Status: types.Null},
			dst:      &float64Slice,
			expected: ([]float64)(nil),
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
		src types.Float8Array
		dst interface{}
	}{
		{
			src: types.Float8Array{
				Elements:   []types.Float8{{Status: types.Null}},
				Dimensions: []types.ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     types.Present,
			},
			dst: &float64Slice,
		},
	}

	for i, tt := range errorTests {
		err := tt.src.AssignTo(tt.dst)
		if err == nil {
			t.Errorf("%d: expected error but none was returned (%v -> %v)", i, tt.src, tt.dst)
		}
	}

}
