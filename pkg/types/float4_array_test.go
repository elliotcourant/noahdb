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

func TestFloat4ArrayTranscode(t *testing.T) {
	testutil.TestSuccessfulTranscode(t, "float4[]", []interface{}{
		&types.Float4Array{
			Elements:   nil,
			Dimensions: nil,
			Status:     types.Present,
		},
		&types.Float4Array{
			Elements: []types.Float4{
				{Float: 1, Status: types.Present},
				{Status: types.Null},
			},
			Dimensions: []types.ArrayDimension{{Length: 2, LowerBound: 1}},
			Status:     types.Present,
		},
		&types.Float4Array{Status: types.Null},
		&types.Float4Array{
			Elements: []types.Float4{
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
		&types.Float4Array{
			Elements: []types.Float4{
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

func TestFloat4ArraySet(t *testing.T) {
	successfulTests := []struct {
		source interface{}
		result types.Float4Array
	}{
		{
			source: []float32{1},
			result: types.Float4Array{
				Elements:   []types.Float4{{Float: 1, Status: types.Present}},
				Dimensions: []types.ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     types.Present},
		},
		{
			source: ([]float32)(nil),
			result: types.Float4Array{Status: types.Null},
		},
	}

	for i, tt := range successfulTests {
		var r types.Float4Array
		err := r.Set(tt.source)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if !reflect.DeepEqual(r, tt.result) {
			t.Errorf("%d: expected %v to convert to %v, but it was %v", i, tt.source, tt.result, r)
		}
	}
}

func TestFloat4ArrayAssignTo(t *testing.T) {
	var float32Slice []float32
	var namedFloat32Slice _float32Slice

	simpleTests := []struct {
		src      types.Float4Array
		dst      interface{}
		expected interface{}
	}{
		{
			src: types.Float4Array{
				Elements:   []types.Float4{{Float: 1.23, Status: types.Present}},
				Dimensions: []types.ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     types.Present,
			},
			dst:      &float32Slice,
			expected: []float32{1.23},
		},
		{
			src: types.Float4Array{
				Elements:   []types.Float4{{Float: 1.23, Status: types.Present}},
				Dimensions: []types.ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     types.Present,
			},
			dst:      &namedFloat32Slice,
			expected: _float32Slice{1.23},
		},
		{
			src:      types.Float4Array{Status: types.Null},
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
		src types.Float4Array
		dst interface{}
	}{
		{
			src: types.Float4Array{
				Elements:   []types.Float4{{Status: types.Null}},
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
