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
	tableTestQueries = []struct {
		Query  string
		Tables []string
	}{
		{
			Query:  "SELECT $1::text;",
			Tables: []string{},
		},
		{
			Query:  "SELECT e.typdelim FROM pg_catalog.pg_type t, pg_catalog.pg_type e WHERE t.oid = $1 and t.typelem = e.oid",
			Tables: []string{"pg_type"},
		},
		{
			Query:  "SELECT e.typdelim FROM pg_catalog.pg_type t, pg_catalog.pg_type e WHERE t.oid = $1 and t.typelem = e.oid AND $2=$3",
			Tables: []string{"pg_type"},
		},
		{
			Query:  "SELECT e.typdelim FROM pg_catalog.pg_type t, pg_catalog.pg_type e WHERE t.oid = $1 and t.typelem = e.oid AND $2=$1",
			Tables: []string{"pg_type"},
		},
		{
			Query:  "SELECT products.id FROM products JOIN types ON types.id=products.type_id",
			Tables: []string{"products", "types"},
		},
		{
			Query:  "SELECT products.id FROM products JOIN types ON types.id=products.type_id WHERE products.id IN (SELECT id FROM other)",
			Tables: []string{"products", "types", "other"},
		},
		{
			Query:  "INSERT INTO products (id) VALUES(1);",
			Tables: []string{"products"},
		},
		{
			Query:  "UPDATE variations SET id=4 WHERE id=3;",
			Tables: []string{"variations"},
		},
	}
)

func Test_GetTables(t *testing.T) {
	for _, item := range tableTestQueries {
		parsed, err := ast.Parse(item.Query)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		stmt := parsed.Statements[0].(ast.RawStmt).Stmt

		tableCount := GetTables(stmt)

		assert.Equal(t, item.Tables, tableCount, "number of tables does not match expected")
	}
}
