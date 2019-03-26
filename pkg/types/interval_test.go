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

func TestIntervalTranscode(t *testing.T) {
	testutil.TestSuccessfulTranscode(t, "interval", []interface{}{
		&types.Interval{Microseconds: 1, Status: types.Present},
		&types.Interval{Microseconds: 1000000, Status: types.Present},
		&types.Interval{Microseconds: 1000001, Status: types.Present},
		&types.Interval{Microseconds: 123202800000000, Status: types.Present},
		&types.Interval{Days: 1, Status: types.Present},
		&types.Interval{Months: 1, Status: types.Present},
		&types.Interval{Months: 12, Status: types.Present},
		&types.Interval{Months: 13, Days: 15, Microseconds: 1000001, Status: types.Present},
		&types.Interval{Microseconds: -1, Status: types.Present},
		&types.Interval{Microseconds: -1000000, Status: types.Present},
		&types.Interval{Microseconds: -1000001, Status: types.Present},
		&types.Interval{Microseconds: -123202800000000, Status: types.Present},
		&types.Interval{Days: -1, Status: types.Present},
		&types.Interval{Months: -1, Status: types.Present},
		&types.Interval{Months: -12, Status: types.Present},
		&types.Interval{Months: -13, Days: -15, Microseconds: -1000001, Status: types.Present},
		&types.Interval{Status: types.Null},
	})
}

func TestIntervalNormalize(t *testing.T) {
	testutil.TestSuccessfulNormalize(t, []testutil.NormalizeTest{
		{
			SQL:   "select '1 second'::interval",
			Value: &types.Interval{Microseconds: 1000000, Status: types.Present},
		},
		{
			SQL:   "select '1.000001 second'::interval",
			Value: &types.Interval{Microseconds: 1000001, Status: types.Present},
		},
		{
			SQL:   "select '34223 hours'::interval",
			Value: &types.Interval{Microseconds: 123202800000000, Status: types.Present},
		},
		{
			SQL:   "select '1 day'::interval",
			Value: &types.Interval{Days: 1, Status: types.Present},
		},
		{
			SQL:   "select '1 month'::interval",
			Value: &types.Interval{Months: 1, Status: types.Present},
		},
		{
			SQL:   "select '1 year'::interval",
			Value: &types.Interval{Months: 12, Status: types.Present},
		},
		{
			SQL:   "select '-13 mon'::interval",
			Value: &types.Interval{Months: -13, Status: types.Present},
		},
	})
}
