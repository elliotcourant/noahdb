package ast

import (
	"testing"
)

func Test_DropStmt_Table(t *testing.T) {
	DoTest(t, DeparseTest{
		Query:    `DROP TABLE temp;`,
		Expected: `DROP TABLE "temp"`,
	})
	DoTest(t, DeparseTest{
		Query:    `drop table temp;`,
		Expected: `DROP TABLE "temp"`,
	})
	DoTest(t, DeparseTest{
		Query:    `drop table "temp"`,
		Expected: `DROP TABLE "temp"`,
	})
}

func Test_DropStmt_Table_Cascade(t *testing.T) {
	DoTest(t, DeparseTest{
		Query:    `DROP TABLE temp CASCADE;`,
		Expected: `DROP TABLE "temp" CASCADE`,
	})
	DoTest(t, DeparseTest{
		Query:    `drop table temp cascade;`,
		Expected: `DROP TABLE "temp" CASCADE`,
	})
	DoTest(t, DeparseTest{
		Query:    `drop table "temp" cascade`,
		Expected: `DROP TABLE "temp" CASCADE`,
	})
}

func Test_DropStmt_Table_IfExists(t *testing.T) {
	DoTest(t, DeparseTest{
		Query:    `DROP TABLE IF EXISTS temp;`,
		Expected: `DROP TABLE IF EXISTS "temp"`,
	})
	DoTest(t, DeparseTest{
		Query:    `drop table if exists temp;`,
		Expected: `DROP TABLE IF EXISTS "temp"`,
	})
	DoTest(t, DeparseTest{
		Query:    `drop table if exists "temp" CASCADE`,
		Expected: `DROP TABLE IF EXISTS "temp" CASCADE`,
	})
}

func Test_DropStmt_AccessMethod(t *testing.T) {
	DoTest(t, DeparseTest{
		Query:    `DROP ACCESS METHOD temp;`,
		Expected: `DROP ACCESS METHOD "temp"`,
	})
}

func Test_DropStmt_Aggregate(t *testing.T) {
	DoTest(t, DeparseTest{
		Query:    `DROP AGGREGATE temp(bigint);`,
		Expected: `DROP AGGREGATE temp(bigint)`,
	})
	DoTest(t, DeparseTest{
		Query:    `DROP AGGREGATE public.temp(integer);`,
		Expected: `DROP AGGREGATE public.temp(int)`,
	})
	DoTest(t, DeparseTest{
		Query:    `DROP AGGREGATE public.temp(integer, text);`,
		Expected: `DROP AGGREGATE public.temp(int, text)`,
	})
}

func Test_DropStmt_Cast(t *testing.T) {
	DoTest(t, DeparseTest{
		Query:    `DROP CAST (text AS int);`,
		Expected: `DROP CAST (text AS int)`,
	})
	DoTest(t, DeparseTest{
		Query:    `DROP CAST (text AS integer);`,
		Expected: `DROP CAST (text AS int)`,
	})
	DoTest(t, DeparseTest{
		Query:    `DROP CAST (BOOL AS text);`,
		Expected: `DROP CAST (bool AS text)`,
	})
}

func Test_DropStmt_Collation(t *testing.T) {
	DoTest(t, DeparseTest{
		Query:    `DROP COLLATION thing;`,
		Expected: `DROP COLLATION "thing"`,
	})
	DoTest(t, DeparseTest{
		Query:    `DROP COLLATION THING;`,
		Expected: `DROP COLLATION "thing"`,
	})
	DoTest(t, DeparseTest{
		Query:    `DROP COLLATION "thing";`,
		Expected: `DROP COLLATION "thing"`,
	})
}

func Test_DropStmt_Conversions(t *testing.T) {
	DoTest(t, DeparseTest{
		Query:    `DROP CONVERSION thing;`,
		Expected: `DROP CONVERSION "thing"`,
	})
	DoTest(t, DeparseTest{
		Query:    `DROP CONVERSION THING;`,
		Expected: `DROP CONVERSION "thing"`,
	})
	DoTest(t, DeparseTest{
		Query:    `DROP CONVERSION "thing";`,
		Expected: `DROP CONVERSION "thing"`,
	})
}

func Test_DropStmt_Database(t *testing.T) {
	DoTest(t, DeparseTest{
		Query:    `DROP DATABASE thing;`,
		Expected: `DROP DATABASE thing`,
	})
	DoTest(t, DeparseTest{
		Query:    `drop database thing;`,
		Expected: `DROP DATABASE thing`,
	})
	DoTest(t, DeparseTest{
		Query:    `drop database IF EXISTS thing;`,
		Expected: `DROP DATABASE IF EXISTS thing`,
	})
}
