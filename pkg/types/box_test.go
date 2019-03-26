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

func TestBoxTranscode(t *testing.T) {
	testutil.TestSuccessfulTranscode(t, "box", []interface{}{
		&types.Box{
			P:      [2]types.Vec2{{7.1, 5.2345678}, {3.14, 1.678}},
			Status: types.Present,
		},
		&types.Box{
			P:      [2]types.Vec2{{7.1, 1.678}, {-13.14, -5.234}},
			Status: types.Present,
		},
		&types.Box{Status: types.Null},
	})
}

func TestBoxNormalize(t *testing.T) {
	testutil.TestSuccessfulNormalize(t, []testutil.NormalizeTest{
		{
			SQL: "select '3.14, 1.678, 7.1, 5.234'::box",
			Value: &types.Box{
				P:      [2]types.Vec2{{7.1, 5.234}, {3.14, 1.678}},
				Status: types.Present,
			},
		},
	})
}
