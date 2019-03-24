package pg_query

import (
	"testing"
)

func Test_TransactionStmt_Generic(t *testing.T) {
	DoTest(t, DeparseTest{
		Query:    `begin`,
		Expected: `BEGIN`,
	})
	DoTest(t, DeparseTest{
		Query:    `commit`,
		Expected: `COMMIT`,
	})
	DoTest(t, DeparseTest{
		Query:    `rollback`,
		Expected: `ROLLBACK`,
	})
}

func Test_TransactionStmt_TwoPhaseCommit(t *testing.T) {
	DoTest(t, DeparseTest{
		Query:    `prepare transaction '1234'`,
		Expected: `PREPARE TRANSACTION '1234'`,
	})
	DoTest(t, DeparseTest{
		Query:    `rollback prepared '1234'`,
		Expected: `ROLLBACK PREPARED '1234'`,
	})
	DoTest(t, DeparseTest{
		Query:    `commit prepared '1234'`,
		Expected: `COMMIT PREPARED '1234'`,
	})
}
