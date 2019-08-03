package datagrip

import (
	"database/sql"
	"github.com/elliotcourant/noahdb/testutils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDataGripRegression(t *testing.T) {
	colony, cleanup := testutils.NewPgTestColony(t)
	defer cleanup()

	t.Run("version", func(t *testing.T) {
		db, err := sql.Open("postgres", testutils.ConnectionString(colony.Addr()))
		if err != nil {
			panic(err)
		}
		defer db.Close()

		stmt, err := db.Prepare(`SELECT version()`)
		assert.NoError(t, err)
		row := stmt.QueryRow()
		var version string
		if err := row.Scan(&version); !assert.NoError(t, err) {
			panic(err)
		}
		assert.NotEmpty(t, version)
	})

	t.Run("keep alive", func(t *testing.T) {
		db, err := sql.Open("postgres", testutils.ConnectionString(colony.Addr()))
		if err != nil {
			panic(err)
		}
		defer db.Close()

		stmt, err := db.Prepare(`SELECT 'keep alive'`)
		assert.NoError(t, err)
		row := stmt.QueryRow()
		var keepAlive string
		if err := row.Scan(&keepAlive); !assert.NoError(t, err) {
			panic(err)
		}
		assert.NotEmpty(t, keepAlive)
		assert.Equal(t, "keep alive", keepAlive)
	})

	t.Run("current database and schemas", func(t *testing.T) {
		db, err := sql.Open("postgres", testutils.ConnectionString(colony.Addr()))
		if err != nil {
			panic(err)
		}
		defer db.Close()

		stmt, err := db.Prepare(`select current_database() as a, current_schemas(false) as b`)
		assert.NoError(t, err)
		row := stmt.QueryRow()
		currentDatabase, currentSchemas := "", ""
		if err := row.Scan(&currentDatabase, &currentSchemas); !assert.NoError(t, err) {
			panic(err)
		}
		assert.NotEmpty(t, currentDatabase)
		assert.NotEmpty(t, currentSchemas)
	})
}
