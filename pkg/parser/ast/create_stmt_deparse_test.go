package pg_query

import (
	"testing"
)

func Test_CreateStmt_Generic(t *testing.T) {
	DoTest(t, DeparseTest{
		Query:    `CREATE TABLE test (id BIGSERIAL PRIMARY KEY, name TEXT);`,
		Expected: `CREATE TABLE "test" (id bigserial PRIMARY KEY, name text)`,
	})
}

func Test_CreateStmt_Tablespace(t *testing.T) {
	DoTest(t, DeparseTest{
		Query:    `CREATE TABLE test (id BIGSERIAL PRIMARY KEY, name TEXT) TABLESPACE thing;`,
		Expected: `CREATE TABLE "test" (id bigserial PRIMARY KEY, name text) TABLESPACE "thing"`,
	})
}

func Test_CreateStmt_ReferenceColumn(t *testing.T) {
	DoTest(t, DeparseTest{
		Query:    `CREATE TABLE public.users (user_id BIGSERIAL PRIMARY KEY, account_id BIGINT NOT NULL REFERENCES public.accounts (account_id), user_number BIGINT);`,
		Expected: `CREATE TABLE "public"."users" (user_id bigserial PRIMARY KEY, account_id bigint NOT NULL REFERENCES "public"."accounts" ("account_id"), user_number bigint)`,
	})
}
