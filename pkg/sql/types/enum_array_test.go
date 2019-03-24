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

func TestEnumArrayTranscode(t *testing.T) {
	setupConn := testutil.MustConnectPgx(t)
	defer testutil.MustClose(t, setupConn)

	if _, err := setupConn.Exec("drop type if exists color"); err != nil {
		t.Fatal(err)
	}
	if _, err := setupConn.Exec("create type color as enum ('red', 'green', 'blue')"); err != nil {
		t.Fatal(err)
	}

	testutil.TestSuccessfulTranscode(t, "color[]", []interface{}{
		&types.EnumArray{
			Elements:   nil,
			Dimensions: nil,
			Status:     types.Present,
		},
		&types.EnumArray{
			Elements: []types.GenericText{
				{String: "red", Status: types.Present},
				{Status: types.Null},
			},
			Dimensions: []types.ArrayDimension{{Length: 2, LowerBound: 1}},
			Status:     types.Present,
		},
		&types.EnumArray{Status: types.Null},
		&types.EnumArray{
			Elements: []types.GenericText{
				{String: "red", Status: types.Present},
				{String: "green", Status: types.Present},
				{String: "blue", Status: types.Present},
				{String: "red", Status: types.Present},
			},
			Dimensions: []types.ArrayDimension{
				{Length: 2, LowerBound: 4},
				{Length: 2, LowerBound: 2},
			},
			Status: types.Present,
		},
	})
}

func TestEnumArrayArraySet(t *testing.T) {
	successfulTests := []struct {
		source interface{}
		result types.EnumArray
	}{
		{
			source: []string{"foo"},
			result: types.EnumArray{
				Elements:   []types.GenericText{{String: "foo", Status: types.Present}},
				Dimensions: []types.ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     types.Present},
		},
		{
			source: ([]string)(nil),
			result: types.EnumArray{Status: types.Null},
		},
	}

	for i, tt := range successfulTests {
		var r types.EnumArray
		err := r.Set(tt.source)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if !reflect.DeepEqual(r, tt.result) {
			t.Errorf("%d: expected %v to convert to %v, but it was %v", i, tt.source, tt.result, r)
		}
	}
}

func TestEnumArrayArrayAssignTo(t *testing.T) {
	var stringSlice []string
	type _stringSlice []string
	var namedStringSlice _stringSlice

	simpleTests := []struct {
		src      types.EnumArray
		dst      interface{}
		expected interface{}
	}{
		{
			src: types.EnumArray{
				Elements:   []types.GenericText{{String: "foo", Status: types.Present}},
				Dimensions: []types.ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     types.Present,
			},
			dst:      &stringSlice,
			expected: []string{"foo"},
		},
		{
			src: types.EnumArray{
				Elements:   []types.GenericText{{String: "bar", Status: types.Present}},
				Dimensions: []types.ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     types.Present,
			},
			dst:      &namedStringSlice,
			expected: _stringSlice{"bar"},
		},
		{
			src:      types.EnumArray{Status: types.Null},
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
		src types.EnumArray
		dst interface{}
	}{
		{
			src: types.EnumArray{
				Elements:   []types.GenericText{{Status: types.Null}},
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
