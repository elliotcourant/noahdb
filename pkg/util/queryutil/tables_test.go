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
			Query:  "SELECT e.typdelim FROM pg_catalog.pg_type t, pg_catalog.pg_type as e WHERE t.oid = $1 and t.typelem = e.oid",
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
			Query:  "SELECT products.id FROM products, things, stuff JOIN types ON types.id=products.type_id",
			Tables: []string{"products", "things", "stuff", "types"},
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

var (
	tableExtendedTestQueries = []struct {
		Query  string
		Tables map[string]string
	}{
		{
			Query:  "SELECT $1::text;",
			Tables: map[string]string{},
		},
		{
			Query: "SELECT e.typdelim FROM pg_catalog.pg_type t, pg_catalog.pg_type as e WHERE t.oid = $1 and t.typelem = e.oid",
			Tables: map[string]string{
				"t": "pg_type",
				"e": "pg_type",
			},
		},
		{
			Query: "SELECT e.typdelim FROM pg_catalog.pg_type t, pg_catalog.pg_type e WHERE t.oid = $1 and t.typelem = e.oid AND $2=$3",
			Tables: map[string]string{
				"t": "pg_type",
				"e": "pg_type",
			},
		},
		{
			Query: "SELECT e.typdelim FROM pg_catalog.pg_type t, pg_catalog.pg_type e WHERE t.oid = $1 and t.typelem = e.oid AND $2=$1",
			Tables: map[string]string{
				"t": "pg_type",
				"e": "pg_type",
			},
		},
		{
			Query: "SELECT products.id FROM products JOIN types ON types.id=products.type_id",
			Tables: map[string]string{
				"products": "products",
				"types":    "types",
			},
		},
		{
			Query: "SELECT products.id FROM products JOIN types ON types.id=products.type_id WHERE products.id IN (SELECT id FROM other)",
			Tables: map[string]string{
				"products": "products",
				"types":    "types",
				"other":    "other",
			},
		},
		{
			Query: "INSERT INTO products (id) VALUES(1);",
			Tables: map[string]string{
				"products": "products",
			},
		},
		{
			Query: "UPDATE variations SET id=4 WHERE id=3;",
			Tables: map[string]string{
				"variations": "variations",
			},
		},
	}
)

func Test_GetExtendedTables(t *testing.T) {
	for _, item := range tableExtendedTestQueries {
		t.Run(item.Query, func(t *testing.T) {
			parsed, err := ast.Parse(item.Query)
			if err != nil {
				t.Error(err)
				t.FailNow()
			}

			stmt := parsed.Statements[0].(ast.RawStmt).Stmt

			tables := GetExtendedTables(stmt)

			assert.Equal(t, item.Tables, tables, "resulting tables do not match expected")
		})
	}
}
