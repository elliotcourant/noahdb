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
	"testing"

	"github.com/elliotcourant/noahdb/pkg/types"
	"github.com/elliotcourant/noahdb/pkg/types/testutil"
)

func TestBPCharArrayTranscode(t *testing.T) {
	testutil.TestSuccessfulTranscode(t, "char(8)[]", []interface{}{
		&types.BPCharArray{
			Elements:   nil,
			Dimensions: nil,
			Status:     types.Present,
		},
		&types.BPCharArray{
			Elements: []types.BPChar{
				{String: "foo     ", Status: types.Present},
				{Status: types.Null},
			},
			Dimensions: []types.ArrayDimension{{Length: 2, LowerBound: 1}},
			Status:     types.Present,
		},
		&types.BPCharArray{Status: types.Null},
		&types.BPCharArray{
			Elements: []types.BPChar{
				{String: "bar     ", Status: types.Present},
				{String: "NuLL    ", Status: types.Present},
				{String: `wow"quz\`, Status: types.Present},
				{String: "1       ", Status: types.Present},
				{String: "1       ", Status: types.Present},
				{String: "null    ", Status: types.Present},
			},
			Dimensions: []types.ArrayDimension{
				{Length: 3, LowerBound: 1},
				{Length: 2, LowerBound: 1},
			},
			Status: types.Present,
		},
		&types.BPCharArray{
			Elements: []types.BPChar{
				{String: " bar    ", Status: types.Present},
				{String: "    baz ", Status: types.Present},
				{String: "    quz ", Status: types.Present},
				{String: "foo     ", Status: types.Present},
			},
			Dimensions: []types.ArrayDimension{
				{Length: 2, LowerBound: 4},
				{Length: 2, LowerBound: 2},
			},
			Status: types.Present,
		},
	})
}
