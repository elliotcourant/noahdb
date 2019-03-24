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

	"github.com/elliotcourant/noahdb/pkg/sql/types"
	"github.com/elliotcourant/noahdb/pkg/sql/types/testutil"
)

func TestPathTranscode(t *testing.T) {
	testutil.TestSuccessfulTranscode(t, "path", []interface{}{
		&types.Path{
			P:      []types.Vec2{{3.14, 1.678901234}, {7.1, 5.234}},
			Closed: false,
			Status: types.Present,
		},
		&types.Path{
			P:      []types.Vec2{{3.14, 1.678}, {7.1, 5.234}, {23.1, 9.34}},
			Closed: true,
			Status: types.Present,
		},
		&types.Path{
			P:      []types.Vec2{{7.1, 1.678}, {-13.14, -5.234}},
			Closed: true,
			Status: types.Present,
		},
		&types.Path{Status: types.Null},
	})
}
