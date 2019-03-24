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

package pg_query

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
