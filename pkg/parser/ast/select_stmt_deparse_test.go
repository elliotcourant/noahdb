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

func Test_SelectStmt_Simple(t *testing.T) {
	DoTest(t, DeparseTest{
		Query:    `SELECT 1;`,
		Expected: `SELECT 1`,
	})
}

func Test_SelectStmt_OrderBy(t *testing.T) {
	DoTest(t, DeparseTest{
		Query:    `SELECT id FROM users ORDER BY id DESC;`,
		Expected: `SELECT "id" FROM "users" ORDER BY "id" DESC`,
	})
}

func Test_SelectStmt_WeirdBoolExpr(t *testing.T) {
	DoTest(t, DeparseTest{
		Query: `select N.oid::bigint as id, datname as name, D.description
from pg_catalog.pg_database N
  left join pg_catalog.pg_shdescription D on N.oid = D.objoid
where not datistemplate
order by case when datname = pg_catalog.current_database() then -1::bigint else N.oid::bigint end`,
		Expected: `SELECT "n"."oid"::bigint AS id, "datname" AS name, "d"."description" FROM "pg_catalog"."pg_database" n LEFT JOIN "pg_catalog"."pg_shdescription" d ON "n"."oid" = "d"."objoid" WHERE NOT datistemplate ORDER BY CASE WHEN "datname" = current_database() THEN 1::bigint ELSE "n"."oid"::bigint END`,
	})
}

func Test_SelectStmt_DataGrip(t *testing.T) {
	DoTest(t, DeparseTest{
		Query: `select t.oid,
						case when nsp.nspname in ('pg_catalog', 'public') then t.typname
							else nsp.nspname||'.'||t.typname
						end
					from pg_type t
					left join pg_type base_type on t.typelem=base_type.oid
					left join pg_namespace nsp on t.typnamespace=nsp.oid
					where (
						  (t.typtype in('b', 'p', 'r', 'e') or 1=1)
						  and (base_type.oid is null or base_type.typtype in('b', 'p', 'r'))
						);`,
		Expected: `SELECT "t"."oid", CASE WHEN "nsp"."nspname" IN ('pg_catalog', 'public') THEN "t"."typname" ELSE "nsp"."nspname" || '.' || "t"."typname" END FROM "pg_type" t LEFT JOIN "pg_type" base_type ON "t"."typelem" = "base_type"."oid" LEFT JOIN "pg_namespace" nsp ON "t"."typnamespace" = "nsp"."oid" WHERE ("t"."typtype" IN ('b', 'p', 'r', 'e') OR 1 = 1) AND ("base_type"."oid" IS NULL OR "base_type"."typtype" IN ('b', 'p', 'r'))`,
	})
}

func Test_SelectStmt_CurrentTimestamp(t *testing.T) {
	DoTest(t, DeparseTest{
		Query:    `select current_timestamp`,
		Expected: `SELECT CURRENT_TIMESTAMP`,
	})
}

func Test_SelectStmt_FunctionCall(t *testing.T) {
	DoTest(t, DeparseTest{
		Query:    `select current_database() as a, current_schemas(false) as b, totalRecords() as c`,
		Expected: `SELECT pg_catalog.current_database() AS a, pg_catalog.current_schemas(false) AS b, pg_catalog.totalrecords() AS c`,
	})
}

func Test_SelectStmt_Between(t *testing.T) {
	DoTest(t, DeparseTest{
		Query:    `select * from test where id between 1 and 3`,
		Expected: `SELECT * FROM "test" WHERE "id" BETWEEN 1 AND 3`,
	})
}

func Test_SelectStmt_NullIf(t *testing.T) {
	DoTest(t, DeparseTest{
		Query:    `select nullif(id,1) from test`,
		Expected: `SELECT NULLIF("id", 1) FROM "test"`,
	})
}

func Test_SelectStmt_Weird_Params(t *testing.T) {
	DoTest(t, DeparseTest{
		Query:    `SELECT n.nspname = ANY(current_schemas(true)), n.nspname, t.typname FROM pg_catalog.pg_type t JOIN pg_catalog.pg_namespace n ON t.typnamespace = n.oid WHERE t.oid = $1`,
		Expected: `SELECT "n"."nspname"=ANY(pg_catalog.current_schemas(true)), "n"."nspname", "t"."typname" FROM "pg_catalog"."pg_type" t JOIN "pg_catalog"."pg_namespace" n ON "t"."typnamespace" = "n"."oid" WHERE "t"."oid" = $1`,
	})
}
