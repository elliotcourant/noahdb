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
	"net"
	"reflect"
	"testing"

	"github.com/elliotcourant/noahdb/pkg/types"
	"github.com/elliotcourant/noahdb/pkg/types/testutil"
)

func TestInetArrayTranscode(t *testing.T) {
	testutil.TestSuccessfulTranscode(t, "inet[]", []interface{}{
		&types.InetArray{
			Elements:   nil,
			Dimensions: nil,
			Status:     types.Present,
		},
		&types.InetArray{
			Elements: []types.Inet{
				{IPNet: mustParseCIDR(t, "12.34.56.0/32"), Status: types.Present},
				{Status: types.Null},
			},
			Dimensions: []types.ArrayDimension{{Length: 2, LowerBound: 1}},
			Status:     types.Present,
		},
		&types.InetArray{Status: types.Null},
		&types.InetArray{
			Elements: []types.Inet{
				{IPNet: mustParseCIDR(t, "127.0.0.1/32"), Status: types.Present},
				{IPNet: mustParseCIDR(t, "12.34.56.0/32"), Status: types.Present},
				{IPNet: mustParseCIDR(t, "192.168.0.1/32"), Status: types.Present},
				{IPNet: mustParseCIDR(t, "2607:f8b0:4009:80b::200e/128"), Status: types.Present},
				{Status: types.Null},
				{IPNet: mustParseCIDR(t, "255.0.0.0/8"), Status: types.Present},
			},
			Dimensions: []types.ArrayDimension{{Length: 3, LowerBound: 1}, {Length: 2, LowerBound: 1}},
			Status:     types.Present,
		},
		&types.InetArray{
			Elements: []types.Inet{
				{IPNet: mustParseCIDR(t, "127.0.0.1/32"), Status: types.Present},
				{IPNet: mustParseCIDR(t, "12.34.56.0/32"), Status: types.Present},
				{IPNet: mustParseCIDR(t, "192.168.0.1/32"), Status: types.Present},
				{IPNet: mustParseCIDR(t, "2607:f8b0:4009:80b::200e/128"), Status: types.Present},
			},
			Dimensions: []types.ArrayDimension{
				{Length: 2, LowerBound: 4},
				{Length: 2, LowerBound: 2},
			},
			Status: types.Present,
		},
	})
}

func TestInetArraySet(t *testing.T) {
	successfulTests := []struct {
		source interface{}
		result types.InetArray
	}{
		{
			source: []*net.IPNet{mustParseCIDR(t, "127.0.0.1/32")},
			result: types.InetArray{
				Elements:   []types.Inet{{IPNet: mustParseCIDR(t, "127.0.0.1/32"), Status: types.Present}},
				Dimensions: []types.ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     types.Present},
		},
		{
			source: ([]*net.IPNet)(nil),
			result: types.InetArray{Status: types.Null},
		},
		{
			source: []net.IP{mustParseCIDR(t, "127.0.0.1/32").IP},
			result: types.InetArray{
				Elements:   []types.Inet{{IPNet: mustParseCIDR(t, "127.0.0.1/32"), Status: types.Present}},
				Dimensions: []types.ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     types.Present},
		},
		{
			source: ([]net.IP)(nil),
			result: types.InetArray{Status: types.Null},
		},
	}

	for i, tt := range successfulTests {
		var r types.InetArray
		err := r.Set(tt.source)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if !reflect.DeepEqual(r, tt.result) {
			t.Errorf("%d: expected %v to convert to %v, but it was %v", i, tt.source, tt.result, r)
		}
	}
}

func TestInetArrayAssignTo(t *testing.T) {
	var ipnetSlice []*net.IPNet
	var ipSlice []net.IP

	simpleTests := []struct {
		src      types.InetArray
		dst      interface{}
		expected interface{}
	}{
		{
			src: types.InetArray{
				Elements:   []types.Inet{{IPNet: mustParseCIDR(t, "127.0.0.1/32"), Status: types.Present}},
				Dimensions: []types.ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     types.Present,
			},
			dst:      &ipnetSlice,
			expected: []*net.IPNet{mustParseCIDR(t, "127.0.0.1/32")},
		},
		{
			src: types.InetArray{
				Elements:   []types.Inet{{Status: types.Null}},
				Dimensions: []types.ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     types.Present,
			},
			dst:      &ipnetSlice,
			expected: []*net.IPNet{nil},
		},
		{
			src: types.InetArray{
				Elements:   []types.Inet{{IPNet: mustParseCIDR(t, "127.0.0.1/32"), Status: types.Present}},
				Dimensions: []types.ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     types.Present,
			},
			dst:      &ipSlice,
			expected: []net.IP{mustParseCIDR(t, "127.0.0.1/32").IP},
		},
		{
			src: types.InetArray{
				Elements:   []types.Inet{{Status: types.Null}},
				Dimensions: []types.ArrayDimension{{LowerBound: 1, Length: 1}},
				Status:     types.Present,
			},
			dst:      &ipSlice,
			expected: []net.IP{nil},
		},
		{
			src:      types.InetArray{Status: types.Null},
			dst:      &ipnetSlice,
			expected: ([]*net.IPNet)(nil),
		},
		{
			src:      types.InetArray{Status: types.Null},
			dst:      &ipSlice,
			expected: ([]net.IP)(nil),
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
}
