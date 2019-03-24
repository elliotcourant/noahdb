package pg_query

import (
	"testing"
)

func Test_InsertStmt_Generic(t *testing.T) {
	DoTest(t, DeparseTest{
		Query:    `insert into users (id, email) values (123, 'test@test.com');`,
		Expected: `INSERT INTO "users" (id, email) VALUES (123, 'test@test.com')`,
	})
	DoTest(t, DeparseTest{
		Query:    `INSERT INTO users (id, email) VALUES (123, 'test@test.com');`,
		Expected: `INSERT INTO "users" (id, email) VALUES (123, 'test@test.com')`,
	})
}

func Test_InsertStmt_Returning(t *testing.T) {
	DoTest(t, DeparseTest{
		Query:    `INSERT INTO users (id, email) VALUES (123, 'test@test.com') RETURNING id;`,
		Expected: `INSERT INTO "users" (id, email) VALUES (123, 'test@test.com') RETURNING "id"`,
	})
	DoTest(t, DeparseTest{
		Query:    `INSERT INTO users (id, email) VALUES (123, 'test@test.com') RETURNING id, email;`,
		Expected: `INSERT INTO "users" (id, email) VALUES (123, 'test@test.com') RETURNING "id", "email"`,
	})
	DoTest(t, DeparseTest{
		Query:    `INSERT INTO users (id, email) VALUES (123, 'test@test.com') RETURNING *;`,
		Expected: `INSERT INTO "users" (id, email) VALUES (123, 'test@test.com') RETURNING *`,
	})
}

func Test_InsertStmt_MultiValue(t *testing.T) {
	DoTest(t, DeparseTest{
		Query:    `INSERT INTO users (id, email) VALUES (123, 'test@test.com'), (321, 'email@email.com');`,
		Expected: `INSERT INTO "users" (id, email) VALUES (123, 'test@test.com'), (321, 'email@email.com')`,
	})
}
