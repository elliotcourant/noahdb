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

func TestNameTranscode(t *testing.T) {
	testutil.TestSuccessfulTranscode(t, "name", []interface{}{
		&types.Name{String: "", Status: types.Present},
		&types.Name{String: "foo", Status: types.Present},
		&types.Name{Status: types.Null},
	})
}

func TestNameSet(t *testing.T) {
	successfulTests := []struct {
		source interface{}
		result types.Name
	}{
		{source: "foo", result: types.Name{String: "foo", Status: types.Present}},
		{source: _string("bar"), result: types.Name{String: "bar", Status: types.Present}},
		{source: (*string)(nil), result: types.Name{Status: types.Null}},
	}

	for i, tt := range successfulTests {
		var d types.Name
		err := d.Set(tt.source)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if d != tt.result {
			t.Errorf("%d: expected %v to convert to %v, but it was %v", i, tt.source, tt.result, d)
		}
	}
}

func TestNameAssignTo(t *testing.T) {
	var s string
	var ps *string

	simpleTests := []struct {
		src      types.Name
		dst      interface{}
		expected interface{}
	}{
		{src: types.Name{String: "foo", Status: types.Present}, dst: &s, expected: "foo"},
		{src: types.Name{Status: types.Null}, dst: &ps, expected: (*string)(nil)},
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
		src      types.Name
		dst      interface{}
		expected interface{}
	}{
		{src: types.Name{String: "foo", Status: types.Present}, dst: &ps, expected: "foo"},
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
		src types.Name
		dst interface{}
	}{
		{src: types.Name{Status: types.Null}, dst: &s},
	}

	for i, tt := range errorTests {
		err := tt.src.AssignTo(tt.dst)
		if err == nil {
			t.Errorf("%d: expected error but none was returned (%v -> %v)", i, tt.src, tt.dst)
		}
	}
}
