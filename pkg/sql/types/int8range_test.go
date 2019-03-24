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

func TestInt8rangeTranscode(t *testing.T) {
	testutil.TestSuccessfulTranscode(t, "Int8range", []interface{}{
		&types.Int8range{LowerType: types.Empty, UpperType: types.Empty, Status: types.Present},
		&types.Int8range{Lower: types.Int8{Int: 1, Status: types.Present}, Upper: types.Int8{Int: 10, Status: types.Present}, LowerType: types.Inclusive, UpperType: types.Exclusive, Status: types.Present},
		&types.Int8range{Lower: types.Int8{Int: -42, Status: types.Present}, Upper: types.Int8{Int: -5, Status: types.Present}, LowerType: types.Inclusive, UpperType: types.Exclusive, Status: types.Present},
		&types.Int8range{Lower: types.Int8{Int: 1, Status: types.Present}, LowerType: types.Inclusive, UpperType: types.Unbounded, Status: types.Present},
		&types.Int8range{Upper: types.Int8{Int: 1, Status: types.Present}, LowerType: types.Unbounded, UpperType: types.Exclusive, Status: types.Present},
		&types.Int8range{Status: types.Null},
	})
}

func TestInt8rangeNormalize(t *testing.T) {
	testutil.TestSuccessfulNormalize(t, []testutil.NormalizeTest{
		{
			SQL:   "select Int8range(1, 10, '(]')",
			Value: types.Int8range{Lower: types.Int8{Int: 2, Status: types.Present}, Upper: types.Int8{Int: 11, Status: types.Present}, LowerType: types.Inclusive, UpperType: types.Exclusive, Status: types.Present},
		},
	})
}
