package ast

import (
	"testing"
)

func Test_DeleteStmt_Generic(t *testing.T) {
	DoTest(t, DeparseTest{
		Query:    `delete from thing;`,
		Expected: `DELETE FROM "thing"`,
	})
	DoTest(t, DeparseTest{
		Query:    `delete from thing where accountId = 123;`,
		Expected: `DELETE FROM "thing" WHERE "accountid" = 123`,
	})
}

func Test_DeleteStmt_Returning(t *testing.T) {
	DoTest(t, DeparseTest{
		Query:    `delete from thing returning *;`,
		Expected: `DELETE FROM "thing" RETURNING *`,
	})
	DoTest(t, DeparseTest{
		Query:    `delete from thing where accountId = 123 returning accountId;`,
		Expected: `DELETE FROM "thing" WHERE "accountid" = 123 RETURNING "accountid"`,
	})
	DoTest(t, DeparseTest{
		Query:    `delete from thing where accountId = 123 returning accountId as account_id;`,
		Expected: `DELETE FROM "thing" WHERE "accountid" = 123 RETURNING "accountid" AS account_id`,
	})
}

func Test_DeleteStmt_Using(t *testing.T) {
	DoTest(t, DeparseTest{
		Query:    `delete from thing using stuff where stuff.id=thing.stuff_id;`,
		Expected: `DELETE FROM "thing" USING "stuff" WHERE "stuff"."id" = "thing"."stuff_id"`,
	})
	DoTest(t, DeparseTest{
		Query:    `delete from thing using stuff, other where stuff.id=thing.stuff_id and other.id=thing.other_id;`,
		Expected: `DELETE FROM "thing" USING "stuff", "other" WHERE "stuff"."id" = "thing"."stuff_id" AND "other"."id" = "thing"."other_id"`,
	})
}

func Test_DeleteStmt_Only(t *testing.T) {
	DoTest(t, DeparseTest{
		Query:    `delete from only thing;`,
		Expected: `DELETE FROM ONLY "thing"`,
	})
}

func Test_DeleteStmt_WhereSelect(t *testing.T) {
	DoTest(t, DeparseTest{
		Query:    `DELETE FROM thing WHERE thing.id = (SELECT id FROM stuff);`,
		Expected: `DELETE FROM "thing" WHERE "thing"."id" = (SELECT "id" FROM "stuff")`,
	})
	DoTest(t, DeparseTest{
		Query:    `DELETE FROM thing WHERE thing.id IN (SELECT id FROM stuff);`,
		Expected: `DELETE FROM "thing" WHERE "thing"."id" IN (SELECT "id" FROM "stuff")`,
	})
	DoTest(t, DeparseTest{
		Query:    `DELETE FROM thing WHERE thing.id NOT IN (SELECT id FROM stuff);`,
		Expected: `DELETE FROM "thing" WHERE NOT "thing"."id" IN (SELECT "id" FROM "stuff")`,
	})
}
