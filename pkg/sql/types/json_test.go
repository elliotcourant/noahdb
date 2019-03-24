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
	"bytes"
	"reflect"
	"testing"

	"github.com/elliotcourant/noahdb/pkg/sql/types"
	"github.com/elliotcourant/noahdb/pkg/sql/types/testutil"
)

func TestJSONTranscode(t *testing.T) {
	testutil.TestSuccessfulTranscode(t, "json", []interface{}{
		&types.JSON{Bytes: []byte("{}"), Status: types.Present},
		&types.JSON{Bytes: []byte("null"), Status: types.Present},
		&types.JSON{Bytes: []byte("42"), Status: types.Present},
		&types.JSON{Bytes: []byte(`"hello"`), Status: types.Present},
		&types.JSON{Status: types.Null},
	})
}

func TestJSONSet(t *testing.T) {
	successfulTests := []struct {
		source interface{}
		result types.JSON
	}{
		{source: "{}", result: types.JSON{Bytes: []byte("{}"), Status: types.Present}},
		{source: []byte("{}"), result: types.JSON{Bytes: []byte("{}"), Status: types.Present}},
		{source: ([]byte)(nil), result: types.JSON{Status: types.Null}},
		{source: (*string)(nil), result: types.JSON{Status: types.Null}},
		{source: []int{1, 2, 3}, result: types.JSON{Bytes: []byte("[1,2,3]"), Status: types.Present}},
		{source: map[string]interface{}{"foo": "bar"}, result: types.JSON{Bytes: []byte(`{"foo":"bar"}`), Status: types.Present}},
	}

	for i, tt := range successfulTests {
		var d types.JSON
		err := d.Set(tt.source)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if !reflect.DeepEqual(d, tt.result) {
			t.Errorf("%d: expected %v to convert to %v, but it was %v", i, tt.source, tt.result, d)
		}
	}
}

func TestJSONAssignTo(t *testing.T) {
	var s string
	var ps *string
	var b []byte

	rawStringTests := []struct {
		src      types.JSON
		dst      *string
		expected string
	}{
		{src: types.JSON{Bytes: []byte("{}"), Status: types.Present}, dst: &s, expected: "{}"},
	}

	for i, tt := range rawStringTests {
		err := tt.src.AssignTo(tt.dst)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if *tt.dst != tt.expected {
			t.Errorf("%d: expected %v to assign %v, but result was %v", i, tt.src, tt.expected, *tt.dst)
		}
	}

	rawBytesTests := []struct {
		src      types.JSON
		dst      *[]byte
		expected []byte
	}{
		{src: types.JSON{Bytes: []byte("{}"), Status: types.Present}, dst: &b, expected: []byte("{}")},
		{src: types.JSON{Status: types.Null}, dst: &b, expected: ([]byte)(nil)},
	}

	for i, tt := range rawBytesTests {
		err := tt.src.AssignTo(tt.dst)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if bytes.Compare(tt.expected, *tt.dst) != 0 {
			t.Errorf("%d: expected %v to assign %v, but result was %v", i, tt.src, tt.expected, *tt.dst)
		}
	}

	var mapDst map[string]interface{}
	type structDst struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	var strDst structDst

	unmarshalTests := []struct {
		src      types.JSON
		dst      interface{}
		expected interface{}
	}{
		{src: types.JSON{Bytes: []byte(`{"foo":"bar"}`), Status: types.Present}, dst: &mapDst, expected: map[string]interface{}{"foo": "bar"}},
		{src: types.JSON{Bytes: []byte(`{"name":"John","age":42}`), Status: types.Present}, dst: &strDst, expected: structDst{Name: "John", Age: 42}},
	}
	for i, tt := range unmarshalTests {
		err := tt.src.AssignTo(tt.dst)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if dst := reflect.ValueOf(tt.dst).Elem().Interface(); !reflect.DeepEqual(dst, tt.expected) {
			t.Errorf("%d: expected %v to assign %v, but result was %v", i, tt.src, tt.expected, dst)
		}
	}

	pointerAllocTests := []struct {
		src      types.JSON
		dst      **string
		expected *string
	}{
		{src: types.JSON{Status: types.Null}, dst: &ps, expected: (*string)(nil)},
	}

	for i, tt := range pointerAllocTests {
		err := tt.src.AssignTo(tt.dst)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if *tt.dst == tt.expected {
			t.Errorf("%d: expected %v to assign %v, but result was %v", i, tt.src, tt.expected, *tt.dst)
		}
	}
}
