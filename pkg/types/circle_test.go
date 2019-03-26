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

func TestCircleTranscode(t *testing.T) {
	testutil.TestSuccessfulTranscode(t, "circle", []interface{}{
		&types.Circle{P: types.Vec2{1.234, 5.67890123}, R: 3.5, Status: types.Present},
		&types.Circle{P: types.Vec2{-1.234, -5.6789}, R: 12.9, Status: types.Present},
		&types.Circle{Status: types.Null},
	})
}
