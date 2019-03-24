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

func TestLineTranscode(t *testing.T) {
	conn := testutil.MustConnectPgx(t)
	if _, ok := conn.ConnInfo.DataTypeForName("line"); !ok {
		t.Skip("Skipping due to no line type")
	}

	// line may exist but not be usable on 9.3 :(
	var isPG93 bool
	err := conn.QueryRow("select version() ~ '9.3'").Scan(&isPG93)
	if err != nil {
		t.Fatal(err)
	}
	if isPG93 {
		t.Skip("Skipping due to unimplemented line type in PG 9.3")
	}

	testutil.TestSuccessfulTranscode(t, "line", []interface{}{
		&types.Line{
			A: 1.23, B: 4.56, C: 7.89012345,
			Status: types.Present,
		},
		&types.Line{
			A: -1.23, B: -4.56, C: -7.89,
			Status: types.Present,
		},
		&types.Line{Status: types.Null},
	})
}
