package pg_query

import (
	"testing"
)

func Test_CreateForeignServerStmt_Generic(t *testing.T) {
	DoTest(t, DeparseTest{
		Query:    `CREATE SERVER test FOREIGN DATA WRAPPER postgres_fdw OPTIONS (host 'foo', dbname 'foodb', port '5432');`,
		Expected: `CREATE SERVER test FOREIGN DATA WRAPPER postgres_fdw OPTIONS (host 'foo', dbname 'foodb', port '5432')`,
	})
	DoTest(t, DeparseTest{
		Query:    `CREATE SERVER test TYPE 'type' FOREIGN DATA WRAPPER postgres_fdw OPTIONS (host 'foo', dbname 'foodb', port '5432');`,
		Expected: `CREATE SERVER test TYPE 'type' FOREIGN DATA WRAPPER postgres_fdw OPTIONS (host 'foo', dbname 'foodb', port '5432')`,
	})
	DoTest(t, DeparseTest{
		Query:    `CREATE SERVER test TYPE 'type' VERSION '123' FOREIGN DATA WRAPPER postgres_fdw OPTIONS (host 'foo', dbname 'foodb', port '5432');`,
		Expected: `CREATE SERVER test TYPE 'type' VERSION '123' FOREIGN DATA WRAPPER postgres_fdw OPTIONS (host 'foo', dbname 'foodb', port '5432')`,
	})
}
