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

func TestChar3Transcode(t *testing.T) {
	testutil.TestSuccessfulTranscodeEqFunc(t, "char(3)", []interface{}{
		&types.BPChar{String: "a  ", Status: types.Present},
		&types.BPChar{String: " a ", Status: types.Present},
		&types.BPChar{String: "嗨  ", Status: types.Present},
		&types.BPChar{String: "   ", Status: types.Present},
		&types.BPChar{Status: types.Null},
	}, func(aa, bb interface{}) bool {
		a := aa.(types.BPChar)
		b := bb.(types.BPChar)

		return a.Status == b.Status && a.String == b.String
	})
}

func TestBPCharAssignTo(t *testing.T) {
	var (
		str string
		run rune
	)
	simpleTests := []struct {
		src      types.BPChar
		dst      interface{}
		expected interface{}
	}{
		{src: types.BPChar{String: "simple", Status: types.Present}, dst: &str, expected: "simple"},
		{src: types.BPChar{String: "嗨", Status: types.Present}, dst: &run, expected: '嗨'},
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

}
