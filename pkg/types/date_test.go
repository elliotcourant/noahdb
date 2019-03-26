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
	"time"

	"github.com/elliotcourant/noahdb/pkg/types"
	"github.com/elliotcourant/noahdb/pkg/types/testutil"
)

func TestDateTranscode(t *testing.T) {
	testutil.TestSuccessfulTranscodeEqFunc(t, "date", []interface{}{
		&types.Date{Time: time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC), Status: types.Present},
		&types.Date{Time: time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC), Status: types.Present},
		&types.Date{Time: time.Date(1999, 12, 31, 0, 0, 0, 0, time.UTC), Status: types.Present},
		&types.Date{Time: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), Status: types.Present},
		&types.Date{Time: time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC), Status: types.Present},
		&types.Date{Time: time.Date(2200, 1, 1, 0, 0, 0, 0, time.UTC), Status: types.Present},
		&types.Date{Status: types.Null},
		&types.Date{Status: types.Present, InfinityModifier: types.Infinity},
		&types.Date{Status: types.Present, InfinityModifier: -types.Infinity},
	}, func(a, b interface{}) bool {
		at := a.(types.Date)
		bt := b.(types.Date)

		return at.Time.Equal(bt.Time) && at.Status == bt.Status && at.InfinityModifier == bt.InfinityModifier
	})
}

func TestDateSet(t *testing.T) {
	type _time time.Time

	successfulTests := []struct {
		source interface{}
		result types.Date
	}{
		{source: time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC), result: types.Date{Time: time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC), Status: types.Present}},
		{source: time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC), result: types.Date{Time: time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC), Status: types.Present}},
		{source: time.Date(1999, 12, 31, 0, 0, 0, 0, time.UTC), result: types.Date{Time: time.Date(1999, 12, 31, 0, 0, 0, 0, time.UTC), Status: types.Present}},
		{source: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), result: types.Date{Time: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), Status: types.Present}},
		{source: time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC), result: types.Date{Time: time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC), Status: types.Present}},
		{source: time.Date(2200, 1, 1, 0, 0, 0, 0, time.UTC), result: types.Date{Time: time.Date(2200, 1, 1, 0, 0, 0, 0, time.UTC), Status: types.Present}},
		{source: _time(time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)), result: types.Date{Time: time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC), Status: types.Present}},
	}

	for i, tt := range successfulTests {
		var d types.Date
		err := d.Set(tt.source)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if d != tt.result {
			t.Errorf("%d: expected %v to convert to %v, but it was %v", i, tt.source, tt.result, d)
		}
	}
}

func TestDateAssignTo(t *testing.T) {
	var tim time.Time
	var ptim *time.Time

	simpleTests := []struct {
		src      types.Date
		dst      interface{}
		expected interface{}
	}{
		{src: types.Date{Time: time.Date(2015, 1, 1, 0, 0, 0, 0, time.Local), Status: types.Present}, dst: &tim, expected: time.Date(2015, 1, 1, 0, 0, 0, 0, time.Local)},
		{src: types.Date{Time: time.Time{}, Status: types.Null}, dst: &ptim, expected: (*time.Time)(nil)},
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
		src      types.Date
		dst      interface{}
		expected interface{}
	}{
		{src: types.Date{Time: time.Date(2015, 1, 1, 0, 0, 0, 0, time.Local), Status: types.Present}, dst: &ptim, expected: time.Date(2015, 1, 1, 0, 0, 0, 0, time.Local)},
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
		src types.Date
		dst interface{}
	}{
		{src: types.Date{Time: time.Date(2015, 1, 1, 0, 0, 0, 0, time.Local), InfinityModifier: types.Infinity, Status: types.Present}, dst: &tim},
		{src: types.Date{Time: time.Date(2015, 1, 1, 0, 0, 0, 0, time.Local), InfinityModifier: types.NegativeInfinity, Status: types.Present}, dst: &tim},
		{src: types.Date{Time: time.Date(2015, 1, 1, 0, 0, 0, 0, time.Local), Status: types.Null}, dst: &tim},
	}

	for i, tt := range errorTests {
		err := tt.src.AssignTo(tt.dst)
		if err == nil {
			t.Errorf("%d: expected error but none was returned (%v -> %v)", i, tt.src, tt.dst)
		}
	}
}
