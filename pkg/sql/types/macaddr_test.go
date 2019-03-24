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
	"net"
	"reflect"
	"testing"

	"github.com/elliotcourant/noahdb/pkg/sql/types"
	"github.com/elliotcourant/noahdb/pkg/sql/types/testutil"
)

func TestMacaddrTranscode(t *testing.T) {
	testutil.TestSuccessfulTranscode(t, "macaddr", []interface{}{
		&types.Macaddr{Addr: mustParseMacaddr(t, "01:23:45:67:89:ab"), Status: types.Present},
		&types.Macaddr{Status: types.Null},
	})
}

func TestMacaddrSet(t *testing.T) {
	successfulTests := []struct {
		source interface{}
		result types.Macaddr
	}{
		{
			source: mustParseMacaddr(t, "01:23:45:67:89:ab"),
			result: types.Macaddr{Addr: mustParseMacaddr(t, "01:23:45:67:89:ab"), Status: types.Present},
		},
		{
			source: "01:23:45:67:89:ab",
			result: types.Macaddr{Addr: mustParseMacaddr(t, "01:23:45:67:89:ab"), Status: types.Present},
		},
	}

	for i, tt := range successfulTests {
		var r types.Macaddr
		err := r.Set(tt.source)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if !reflect.DeepEqual(r, tt.result) {
			t.Errorf("%d: expected %v to convert to %v, but it was %v", i, tt.source, tt.result, r)
		}
	}
}

func TestMacaddrAssignTo(t *testing.T) {
	{
		src := types.Macaddr{Addr: mustParseMacaddr(t, "01:23:45:67:89:ab"), Status: types.Present}
		var dst net.HardwareAddr
		expected := mustParseMacaddr(t, "01:23:45:67:89:ab")

		err := src.AssignTo(&dst)
		if err != nil {
			t.Error(err)
		}

		if bytes.Compare([]byte(dst), []byte(expected)) != 0 {
			t.Errorf("expected %v to assign %v, but result was %v", src, expected, dst)
		}
	}

	{
		src := types.Macaddr{Addr: mustParseMacaddr(t, "01:23:45:67:89:ab"), Status: types.Present}
		var dst string
		expected := "01:23:45:67:89:ab"

		err := src.AssignTo(&dst)
		if err != nil {
			t.Error(err)
		}

		if dst != expected {
			t.Errorf("expected %v to assign %v, but result was %v", src, expected, dst)
		}
	}
}
