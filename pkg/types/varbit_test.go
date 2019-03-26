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

func TestVarbitTranscode(t *testing.T) {
	testutil.TestSuccessfulTranscode(t, "varbit", []interface{}{
		&types.Varbit{Bytes: []byte{}, Len: 0, Status: types.Present},
		&types.Varbit{Bytes: []byte{0, 1, 128, 254, 255}, Len: 40, Status: types.Present},
		&types.Varbit{Bytes: []byte{0, 1, 128, 254, 128}, Len: 33, Status: types.Present},
		&types.Varbit{Status: types.Null},
	})
}

func TestVarbitNormalize(t *testing.T) {
	t.Skip("cannot transcode from pgx to npgx")
	testutil.TestSuccessfulNormalize(t, []testutil.NormalizeTest{
		{
			SQL:   "select B'111111111'",
			Value: &types.Varbit{Bytes: []byte{255, 128}, Len: 9, Status: types.Present},
		},
	})
}
