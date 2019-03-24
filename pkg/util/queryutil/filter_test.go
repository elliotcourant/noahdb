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

package queryutil

import (
	"github.com/elliotcourant/noahdb/pkg/ast"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_FindAccountIds(t *testing.T) {
	query := `SELECT * FROM users INNER JOIN accounts ON users.account_id=accounts.account_id WHERE users.account_id = 5`
	tree, _ := ast.Parse(query)
	ids := FindAccountIds(tree, "account_id")
	assert.Equal(t, []uint64{5}, ids, "the returned ids do not match expected values")
}
