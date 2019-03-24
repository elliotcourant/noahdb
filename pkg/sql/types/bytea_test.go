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

func TestByteaTranscode(t *testing.T) {
	testutil.TestSuccessfulTranscode(t, "bytea", []interface{}{
		&types.Bytea{Bytes: []byte{1, 2, 3}, Status: types.Present},
		&types.Bytea{Bytes: []byte{}, Status: types.Present},
		&types.Bytea{Bytes: nil, Status: types.Null},
	})
}

func TestByteaSet(t *testing.T) {
	successfulTests := []struct {
		source interface{}
		result types.Bytea
	}{
		{source: []byte{1, 2, 3}, result: types.Bytea{Bytes: []byte{1, 2, 3}, Status: types.Present}},
		{source: []byte{}, result: types.Bytea{Bytes: []byte{}, Status: types.Present}},
		{source: []byte(nil), result: types.Bytea{Status: types.Null}},
		{source: _byteSlice{1, 2, 3}, result: types.Bytea{Bytes: []byte{1, 2, 3}, Status: types.Present}},
		{source: _byteSlice(nil), result: types.Bytea{Status: types.Null}},
	}

	for i, tt := range successfulTests {
		var r types.Bytea
		err := r.Set(tt.source)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if !reflect.DeepEqual(r, tt.result) {
			t.Errorf("%d: expected %v to convert to %v, but it was %v", i, tt.source, tt.result, r)
		}
	}
}

func TestByteaAssignTo(t *testing.T) {
	var buf []byte
	var _buf _byteSlice
	var pbuf *[]byte
	var _pbuf *_byteSlice

	simpleTests := []struct {
		src      types.Bytea
		dst      interface{}
		expected interface{}
	}{
		{src: types.Bytea{Bytes: []byte{1, 2, 3}, Status: types.Present}, dst: &buf, expected: []byte{1, 2, 3}},
		{src: types.Bytea{Bytes: []byte{1, 2, 3}, Status: types.Present}, dst: &_buf, expected: _byteSlice{1, 2, 3}},
		{src: types.Bytea{Bytes: []byte{1, 2, 3}, Status: types.Present}, dst: &pbuf, expected: &[]byte{1, 2, 3}},
		{src: types.Bytea{Bytes: []byte{1, 2, 3}, Status: types.Present}, dst: &_pbuf, expected: &_byteSlice{1, 2, 3}},
		{src: types.Bytea{Status: types.Null}, dst: &pbuf, expected: (*[]byte)(nil)},
		{src: types.Bytea{Status: types.Null}, dst: &_pbuf, expected: (*_byteSlice)(nil)},
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
}
