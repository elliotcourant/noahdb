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
	"time"

	"github.com/elliotcourant/noahdb/pkg/sql/types"
	"github.com/elliotcourant/noahdb/pkg/sql/types/testutil"
)

func TestTsrangeTranscode(t *testing.T) {
	testutil.TestSuccessfulTranscodeEqFunc(t, "tsrange", []interface{}{
		&types.Tsrange{LowerType: types.Empty, UpperType: types.Empty, Status: types.Present},
		&types.Tsrange{
			Lower:     types.Timestamp{Time: time.Date(1990, 12, 31, 0, 0, 0, 0, time.UTC), Status: types.Present},
			Upper:     types.Timestamp{Time: time.Date(2028, 1, 1, 0, 23, 12, 0, time.UTC), Status: types.Present},
			LowerType: types.Inclusive,
			UpperType: types.Exclusive,
			Status:    types.Present,
		},
		&types.Tsrange{
			Lower:     types.Timestamp{Time: time.Date(1800, 12, 31, 0, 0, 0, 0, time.UTC), Status: types.Present},
			Upper:     types.Timestamp{Time: time.Date(2200, 1, 1, 0, 23, 12, 0, time.UTC), Status: types.Present},
			LowerType: types.Inclusive,
			UpperType: types.Exclusive,
			Status:    types.Present,
		},
		&types.Tsrange{Status: types.Null},
	}, func(aa, bb interface{}) bool {
		a := aa.(types.Tsrange)
		b := bb.(types.Tsrange)

		return a.Status == b.Status &&
			a.Lower.Time.Equal(b.Lower.Time) &&
			a.Lower.Status == b.Lower.Status &&
			a.Lower.InfinityModifier == b.Lower.InfinityModifier &&
			a.Upper.Time.Equal(b.Upper.Time) &&
			a.Upper.Status == b.Upper.Status &&
			a.Upper.InfinityModifier == b.Upper.InfinityModifier
	})
}
