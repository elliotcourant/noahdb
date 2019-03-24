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

var (
	testColumnQueries = []struct {
		Query       string
		ColumnCount int
	}{
		{
			Query:       "SELECT $1::text;",
			ColumnCount: 1,
		},
		{
			Query:       "SELECT e.typdelim FROM pg_catalog.pg_type t, pg_catalog.pg_type e WHERE t.oid = $1 and t.typelem = e.oid",
			ColumnCount: 1,
		},
		{
			Query:       "SELECT e.typdelim, e.thing FROM pg_catalog.pg_type t, pg_catalog.pg_type e WHERE t.oid = $1 and t.typelem = e.oid AND $2=$1",
			ColumnCount: 2,
		},
		{
			Query:       "INSERT INTO test (id, thing) VALUES(1, 2);",
			ColumnCount: 0,
		},
		{
			Query:       "INSERT INTO test (id, thing, another_thing) VALUES(1, 2, 3) RETURNING id;",
			ColumnCount: 1,
		},
		{
			Query:       "UPDATE test SET thing = 1 WHERE stuff = 0;",
			ColumnCount: 0,
		},
		{
			Query:       "UPDATE test SET thing = 1, more = true WHERE stuff = 0 RETURNING id;",
			ColumnCount: 1,
		},
		{
			Query:       "DELETE FROM test WHERE stuff = true RETURNING thing;",
			ColumnCount: 1,
		},
		{
			Query:       "DELETE FROM test WHERE stuff = true;",
			ColumnCount: 0,
		},
	}
)

func Test_GetColumns(t *testing.T) {
	for _, item := range testColumnQueries {
		parsed, err := ast.Parse(item.Query)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		stmt := parsed.Statements[0].(ast.RawStmt).Stmt.(ast.Stmt)

		colCount := GetColumns(stmt)

		assert.Equal(t, item.ColumnCount, len(colCount), "number of columns does not match expected")
	}
}
