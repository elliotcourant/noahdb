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

func TestOIDValueTranscode(t *testing.T) {
	testutil.TestSuccessfulTranscode(t, "oid", []interface{}{
		&types.OIDValue{Uint: 42, Status: types.Present},
		&types.OIDValue{Status: types.Null},
	})
}

func TestOIDValueSet(t *testing.T) {
	successfulTests := []struct {
		source interface{}
		result types.OIDValue
	}{
		{source: uint32(1), result: types.OIDValue{Uint: 1, Status: types.Present}},
	}

	for i, tt := range successfulTests {
		var r types.OIDValue
		err := r.Set(tt.source)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if r != tt.result {
			t.Errorf("%d: expected %v to convert to %v, but it was %v", i, tt.source, tt.result, r)
		}
	}
}

func TestOIDValueAssignTo(t *testing.T) {
	var ui32 uint32
	var pui32 *uint32

	simpleTests := []struct {
		src      types.OIDValue
		dst      interface{}
		expected interface{}
	}{
		{src: types.OIDValue{Uint: 42, Status: types.Present}, dst: &ui32, expected: uint32(42)},
		{src: types.OIDValue{Status: types.Null}, dst: &pui32, expected: (*uint32)(nil)},
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
		src      types.OIDValue
		dst      interface{}
		expected interface{}
	}{
		{src: types.OIDValue{Uint: 42, Status: types.Present}, dst: &pui32, expected: uint32(42)},
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
		src types.OIDValue
		dst interface{}
	}{
		{src: types.OIDValue{Status: types.Null}, dst: &ui32},
	}

	for i, tt := range errorTests {
		err := tt.src.AssignTo(tt.dst)
		if err == nil {
			t.Errorf("%d: expected error but none was returned (%v -> %v)", i, tt.src, tt.dst)
		}
	}
}
