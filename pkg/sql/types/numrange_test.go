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
	"math/big"
	"testing"

	"github.com/elliotcourant/noahdb/pkg/sql/types"
	"github.com/elliotcourant/noahdb/pkg/sql/types/testutil"
)

func TestNumrangeTranscode(t *testing.T) {
	testutil.TestSuccessfulTranscode(t, "numrange", []interface{}{
		&types.Numrange{
			LowerType: types.Empty,
			UpperType: types.Empty,
			Status:    types.Present,
		},
		&types.Numrange{
			Lower:     types.Numeric{Int: big.NewInt(-543), Exp: 3, Status: types.Present},
			Upper:     types.Numeric{Int: big.NewInt(342), Exp: 1, Status: types.Present},
			LowerType: types.Inclusive,
			UpperType: types.Exclusive,
			Status:    types.Present,
		},
		&types.Numrange{
			Lower:     types.Numeric{Int: big.NewInt(-42), Exp: 1, Status: types.Present},
			Upper:     types.Numeric{Int: big.NewInt(-5), Exp: 0, Status: types.Present},
			LowerType: types.Inclusive,
			UpperType: types.Exclusive,
			Status:    types.Present,
		},
		&types.Numrange{
			Lower:     types.Numeric{Int: big.NewInt(-42), Exp: 1, Status: types.Present},
			LowerType: types.Inclusive,
			UpperType: types.Unbounded,
			Status:    types.Present,
		},
		&types.Numrange{
			Upper:     types.Numeric{Int: big.NewInt(-42), Exp: 1, Status: types.Present},
			LowerType: types.Unbounded,
			UpperType: types.Exclusive,
			Status:    types.Present,
		},
		&types.Numrange{Status: types.Null},
	})
}
