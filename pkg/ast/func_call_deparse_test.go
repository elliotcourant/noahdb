package ast

import (
	"testing"
)

func Test_FuncCall_Generic(t *testing.T) {
	DoTest(t, DeparseTest{
		Query:    `select current_database() as a, current_schemas(false) as b`,
		Expected: `SELECT pg_catalog.current_database() AS a, pg_catalog.current_schemas(false) AS b`,
	})
}
