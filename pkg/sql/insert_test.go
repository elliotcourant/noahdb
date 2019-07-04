package sql_test

import (
	"database/sql"
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
}
