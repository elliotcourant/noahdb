package pg_query

import (
	"testing"
)

func Test_VariableSetStmt_Generic(t *testing.T) {
	DoTest(t, DeparseTest{
		Query:    `SET extra_float_digits = 3`,
		Expected: `SET extra_float_digits TO 3`,
	})
}
