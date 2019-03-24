package pg_query

import (
	"testing"
)

func Test_UpdateStmt_Generic(t *testing.T) {
	DoTest(t, DeparseTest{
		Query:    `update users set enabled=1`,
		Expected: `UPDATE "users" SET enabled = 1`,
	})
}

func Test_UpdateStmt_Returning(t *testing.T) {
	DoTest(t, DeparseTest{
		Query:    `UPDATE users set is_enabled=true WHERE user_id='2' returning *`,
		Expected: `UPDATE "users" SET is_enabled = true WHERE "user_id" = '2' RETURNING *`,
	})
}

func Test_UpdateStmt_TextArray(t *testing.T) {
	DoTest(t, DeparseTest{
		Query:    `UPDATE users set is_enabled='{apple,cherry apple, avocado}'::text[] WHERE user_id='2' returning *`,
		Expected: `UPDATE "users" SET is_enabled = '{apple,cherry apple, avocado}'::text[] WHERE "user_id" = '2' RETURNING *`,
	})
}

func Test_UpdateStmt_Numeric(t *testing.T) {
	DoTest(t, DeparseTest{
		Query:    `UPDATE users set is_enabled=500.215::numeric(5,18) WHERE user_id='2'`,
		Expected: `UPDATE "users" SET is_enabled = 500.215::numeric(5, 18) WHERE "user_id" = '2'`,
	})
}

func Test_UpdateStmt_Interval(t *testing.T) {
	DoTest(t, DeparseTest{
		Query:    `UPDATE users set is_enabled=interval '2 months ago' WHERE user_id='2' returning *`,
		Expected: `UPDATE "users" SET is_enabled = '2 months ago'::interval WHERE "user_id" = '2' RETURNING *`,
	})
}
