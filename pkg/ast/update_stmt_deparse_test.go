package ast

import (
	"testing"
)

func Test_UpdateStmt_Generic(t *testing.T) {
	DoTest(t, DeparseTest{
		Query:    `update users set enabled=1`,
		Expected: `UPDATE "users" SET enabled = 1`,
	})
}

func Test_UpdateStmt_Generic_Multiple(t *testing.T) {
	DoTest(t, DeparseTest{
		Query:    `update users set enabled=1, disabled=3`,
		Expected: `UPDATE "users" SET enabled = 1, disabled = 3`,
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

func Test_UpdateStmt_From(t *testing.T) {
	DoTest(t, DeparseTest{
		Query: `UPDATE "books" AS "book"
SET "title"      = COALESCE(_data."title", "book"."title"),
    "author_id"  = COALESCE(_data."author_id", "book"."author_id"),
    "editor_id"  = COALESCE(_data."editor_id", "book"."editor_id"),
    "created_at" = COALESCE(_data."created_at", "book"."created_at"),
    "updated_at" = COALESCE(_data."updated_at", "book"."updated_at")
FROM (VALUES (1::BIGINT, 'updated book 1'::text, NULL::bigint, NULL::bigint, NULL::timestamptz,
              NULL::timestamptz)) AS _data("id", "title", "author_id", "editor_id", "created_at", "updated_at")
WHERE "book"."id" = _data."id" RETURNING *`,
		Expected: `UPDATE "books" book SET title = COALESCE("_data"."title", "book"."title"), author_id = COALESCE("_data"."author_id", "book"."author_id"), editor_id = COALESCE("_data"."editor_id", "book"."editor_id"), created_at = COALESCE("_data"."created_at", "book"."created_at"), updated_at = COALESCE("_data"."updated_at", "book"."updated_at") FROM ( VALUES (1::bigint, 'updated book 1'::text, NULL::bigint, NULL::bigint, NULL::timestamptz, NULL::timestamptz) ) AS _data ( "id", "title", "author_id", "editor_id", "created_at", "updated_at" ) WHERE "book"."id" = "_data"."id" RETURNING *`,
	})
}
