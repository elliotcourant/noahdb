package sql_test

import (
	"database/sql"
	"fmt"
	"github.com/elliotcourant/noahdb/testutils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewInsertStatementPlan(t *testing.T) {
	colony, cleanup := testutils.NewPgTestColony(t)
	defer cleanup()

	db, err := sql.Open("postgres", testutils.ConnectionString(colony.Addr()))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	t.Run("create table and insert", func(t *testing.T) {
		_, err = db.Exec(`CREATE TABLE global_table (id BIGSERIAL NOT NULL PRIMARY KEY, name TEXT UNIQUE);`)
		if !assert.NoError(t, err) {
			panic(err)
		}

		table, ok, err := colony.Tables().GetTable("global_table")
		if !assert.NoError(t, err) {
			panic(err)
		}
		assert.True(t, ok)
		assert.NotEmpty(t, table)

		_, err = db.Exec(`INSERT INTO global_table (name) VALUES('test');`)
		if !assert.NoError(t, err) {
			panic(err)
		}

		rows, err := db.Query(`SELECT id, name FROM global_table;`)
		if !assert.NoError(t, err) {
			panic(err)
		}

		type TempRow struct {
			ID   uint64
			Name string
		}

		values := make([]TempRow, 0)

		for rows.Next() {
			if err := rows.Err(); !assert.NoError(t, err) {
				panic(err)
			}

			value := TempRow{}
			if err := rows.Scan(
				&value.ID,
				&value.Name,
			); !assert.NoError(t, err) {
				panic(err)
			}
			values = append(values, value)
		}

		if err := rows.Err(); !assert.NoError(t, err) {
			panic(err)
		}

		assert.NotEmpty(t, values)
	})

	t.Run("different inserts", func(t *testing.T) {
		_, err = db.Exec(`CREATE TABLE different_table (id BIGSERIAL NOT NULL PRIMARY KEY, name TEXT UNIQUE);`)
		if !assert.NoError(t, err) {
			panic(err)
		}

		table, ok, err := colony.Tables().GetTable("different_table")
		if !assert.NoError(t, err) {
			panic(err)
		}
		assert.True(t, ok)
		assert.NotEmpty(t, table)

		t.Run("normal insert", func(t *testing.T) {
			_, err = db.Exec(`INSERT INTO different_table (name) VALUES('test');`)
			if !assert.NoError(t, err) {
				panic(err)
			}
		})

		t.Run("insert with returning clause", func(t *testing.T) {
			rows, err := db.Query(`INSERT INTO different_table (name) VALUES('test another') RETURNING id;`)
			if !assert.NoError(t, err) {
				panic(err)
			}

			values := make([]uint64, 0)

			for rows.Next() {
				if err := rows.Err(); !assert.NoError(t, err) {
					panic(err)
				}

				var value uint64
				if err := rows.Scan(
					&value,
				); !assert.NoError(t, err) {
					panic(err)
				}
				values = append(values, value)
			}

			if err := rows.Err(); !assert.NoError(t, err) {
				panic(err)
			}

			assert.NotEmpty(t, values)
		})

		t.Run("insert with default serial value", func(t *testing.T) {
			// These two tests make sure that the ID column can be found in any position
			// if the column is specified.
			t.Run("with serial being the first column", func(t *testing.T) {
				_, err = db.Exec(`INSERT INTO different_table (id, name) VALUES(DEFAULT, 'with default');`)
				if !assert.NoError(t, err) {
					panic(err)
				}
			})

			t.Run("with serial being a different column", func(t *testing.T) {
				_, err = db.Exec(`INSERT INTO different_table (name, id) VALUES('with default different', DEFAULT);`)
				if !assert.NoError(t, err) {
					panic(err)
				}
			})
		})

		t.Run("fail if serial value provided", func(t *testing.T) {
			_, err = db.Exec(`INSERT INTO different_table (id, name) VALUES(123, 'with default');`)
			assert.EqualError(t, err, "pq: cannot manually set value of serialized column [id]")
		})
	})

	t.Run("create tenants table and insert", func(t *testing.T) {
		_, err = db.Exec(`CREATE TABLE tenants_table (id BIGSERIAL NOT NULL PRIMARY KEY, name TEXT UNIQUE) TABLESPACE "noah.tenants";`)
		if !assert.NoError(t, err) {
			panic(err)
		}

		table, ok, err := colony.Tables().GetTable("tenants_table")
		if !assert.NoError(t, err) {
			panic(err)
		}
		assert.True(t, ok)
		assert.NotEmpty(t, table)

		_, err = db.Exec(`INSERT INTO tenants_table (name) VALUES('test 1');`)
		if !assert.NoError(t, err) {
			panic(err)
		}

		_, err = db.Exec(`INSERT INTO tenants_table (name) VALUES('test 2');`)
		if !assert.NoError(t, err) {
			panic(err)
		}

		rows, err := db.Query(`SELECT id, name FROM tenants_table;`)
		if !assert.NoError(t, err) {
			panic(err)
		}

		type TempRow struct {
			ID   uint64
			Name string
		}

		values := make([]TempRow, 0)

		for rows.Next() {
			if err := rows.Err(); !assert.NoError(t, err) {
				panic(err)
			}

			value := TempRow{}
			if err := rows.Scan(
				&value.ID,
				&value.Name,
			); !assert.NoError(t, err) {
				panic(err)
			}
			values = append(values, value)
		}

		if err := rows.Err(); !assert.NoError(t, err) {
			panic(err)
		}

		assert.NotEmpty(t, values)

		tenants, err := colony.Tenants().GetTenants()
		assert.NoError(t, err)
		assert.NotEmpty(t, tenants)
	})
}

func TestInsertShardedTablePlan(t *testing.T) {
	colony, cleanup := testutils.NewPgTestColony(t)
	defer cleanup()

	db, err := sql.Open("postgres", testutils.ConnectionString(colony.Addr()))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE tenants_table (id BIGSERIAL NOT NULL PRIMARY KEY, name TEXT UNIQUE) TABLESPACE "noah.tenants";`)
	if !assert.NoError(t, err) {
		panic(err)
	}

	table, ok, err := colony.Tables().GetTable("tenants_table")
	if !assert.NoError(t, err) {
		panic(err)
	}
	assert.True(t, ok)
	assert.NotEmpty(t, table)

	numberOfTenantsToCreate := 15
	for i := 0; i < numberOfTenantsToCreate; i++ {
		_, err = db.Exec(fmt.Sprintf(`INSERT INTO tenants_table (name) VALUES('test %d');`, i))
		if !assert.NoError(t, err) {
			panic(err)
		}
	}

	tenants, err := colony.Tenants().GetTenants()
	if !assert.NoError(t, err) {
		panic(err)
	}

	assert.Equal(t, numberOfTenantsToCreate, len(tenants))

	_, err = db.Exec(`CREATE TABLE users_table (id BIGSERIAL NOT NULL PRIMARY KEY, account_id BIGINT NOT NULL REFERENCES tenants_table (id), name TEXT, CONSTRAINT uq_name UNIQUE (account_id, name)) TABLESPACE "noah.sharded";`)
	if !assert.NoError(t, err) {
		panic(err)
	}

	userTable, ok, err := colony.Tables().GetTable("users_table")
	if !assert.NoError(t, err) {
		panic(err)
	}
	assert.True(t, ok)
	assert.NotEmpty(t, userTable)

	numberOfUsersToInsertPerTenant := 15
	for _, tenant := range tenants {
		for i := 0; i < numberOfUsersToInsertPerTenant; i++ {
			_, err = db.Exec(fmt.Sprintf(`INSERT INTO users_table (account_id, name) VALUES(%d, 'user %d');`, tenant.TenantID, i))
			if !assert.NoError(t, err) {
				panic(err)
			}
		}
	}
}
