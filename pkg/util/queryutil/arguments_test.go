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
	"encoding/json"
	"fmt"
	"github.com/elliotcourant/noahdb/pkg/ast"
	"github.com/elliotcourant/noahdb/pkg/types"
	"github.com/readystock/golog"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	testQueries = []struct {
		Query    string
		ArgCount int
	}{
		{
			Query:    "SELECT $1::text;",
			ArgCount: 1,
		},
		{
			Query:    "SELECT e.typdelim FROM pg_catalog.pg_type t, pg_catalog.pg_type e WHERE t.oid = $1 and t.typelem = e.oid",
			ArgCount: 1,
		},
		{
			Query:    "SELECT e.typdelim FROM pg_catalog.pg_type t, pg_catalog.pg_type e WHERE t.oid = $1 and t.typelem = e.oid AND $2=$3",
			ArgCount: 3,
		},
		{
			Query:    "SELECT e.typdelim FROM pg_catalog.pg_type t, pg_catalog.pg_type e WHERE t.oid = $1 and t.typelem = e.oid AND $2=$1",
			ArgCount: 2,
		},
	}
)

func Test_GetArguments(t *testing.T) {
	for _, item := range testQueries {
		parsed, err := ast.Parse(item.Query)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		stmt := parsed.Statements[0].(ast.RawStmt).Stmt

		argCount := GetArguments(stmt)

		assert.Equal(t, item.ArgCount, len(argCount), "number of arguments does not match expected")
	}
}

var (
	testReplacements = []struct {
		Query     string
		ArgCount  int
		Arguments QueryArguments
	}{
		{
			Query:    "SELECT $1",
			ArgCount: 1,
			Arguments: map[string]types.Value{
				"1": &types.Int4{
					Status: types.Present,
					Int:    1,
				},
			},
		},
		{
			Query:    "SELECT products.id FROM products WHERE products.sku=$1 AND products.type=$2",
			ArgCount: 2,
			Arguments: map[string]types.Value{
				"1": &types.Int4{
					Status: types.Present,
					Int:    1,
				},
				"2": &types.Bool{
					Status: types.Present,
					Bool:   true,
				},
			},
		},
		{
			Query:    "UPDATE users SET enabled=$1 WHERE type=$2",
			ArgCount: 2,
			Arguments: map[string]types.Value{
				"1": &types.Int4{
					Status: types.Present,
					Int:    1,
				},
				"2": &types.Bool{
					Status: types.Present,
					Bool:   true,
				},
			},
		},
		{
			Query:    "INSERT INTO users (id, enabled) VALUES($1, $2) RETURNING *;",
			ArgCount: 2,
			Arguments: map[string]types.Value{
				"1": &types.Int4{
					Status: types.Present,
					Int:    1,
				},
				"2": &types.Bool{
					Status: types.Present,
					Bool:   false,
				},
			},
		},
		{
			Query:    "INSERT INTO users (id, enabled, setup, value) VALUES($1, $2, $1, $3) RETURNING *;",
			ArgCount: 3,
			Arguments: map[string]types.Value{
				"1": &types.Int4{
					Status: types.Present,
					Int:    1,
				},
				"2": &types.Bool{
					Status: types.Present,
					Bool:   false,
				},
				"3": func() *types.Numeric {
					float := types.Numeric{}
					float.Set(float64(5.6))
					return &float
				}(),
			},
		},
		{
			Query:     "select current_database(), current_schema(), current_user",
			ArgCount:  0,
			Arguments: map[string]types.Value{},
		},
		{
			Query:    "INSERT INTO users (id, enabled, setup) VALUES($1, $2, $1) RETURNING *;",
			ArgCount: 2,
			Arguments: map[string]types.Value{
				"1": &types.Int4{
					Status: types.Present,
					Int:    1,
				},
				"2": &types.Text{
					Status: types.Present,
					String: "hello world",
				},
			},
		},
		{
			Query:    "INSERT INTO users (id, enabled, setup) VALUES($1, $2, $1) RETURNING *;",
			ArgCount: 2,
			Arguments: map[string]types.Value{
				"1": &types.Float4{
					Status: types.Present,
					Float:  82.3,
				},
				"2": &types.Float8{
					Status: types.Present,
					Float:  1.4,
				},
			},
		},
		{
			Query:    "INSERT INTO users (id, enabled, setup) VALUES($1, $2, $1) RETURNING *;",
			ArgCount: 2,
			Arguments: map[string]types.Value{
				"1": &types.Float4{
					Status: types.Null,
				},
				"2": &types.Float8{
					Status: types.Present,
					Float:  1.4,
				},
			},
		},
		// {
		//     Query:    "DELETE FROM users WHERE user_id = $1;",
		//     ArgCount: 1,
		//     Arguments: map[string]types.Value{
		//         "1": &types.Int8{
		//             Status: types.Present,
		//             Int:    28412931,
		//         },
		//     },
		// },
	}
)

func Test_ReplaceArguments(t *testing.T) {
	for _, item := range testReplacements {
		parsed, err := ast.Parse(item.Query)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		stmt := parsed.Statements[0].(ast.RawStmt).Stmt

		func() {
			defer func() {
				if r := recover(); r != nil {
					golog.Errorf("Replacing arguments in query `%s` has resulted in a panic", item.Query)
					j, _ := json.Marshal(stmt)
					golog.Errorf("Parse Tree: ->")
					golog.Info(string(j))
					golog.Fatal(r)
				}
			}()

			argCount := GetArguments(stmt)

			assert.Equal(t, item.ArgCount, len(argCount), "number of arguments does not match expected")

			// Now we will replace the arguments, and there should be 0 after
			result := ReplaceArguments(stmt, item.Arguments)

			argCount = GetArguments(result)
			assert.Equal(t, 0, len(argCount), "number of arguments should now be 0")

			query, err := result.(ast.Node).Deparse(ast.Context_None)
			if err != nil {
				t.Error(err)
				t.FailNow()
			}
			fmt.Println("Query: ", query)
		}()
	}
}
