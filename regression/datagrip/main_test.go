package datagrip

import (
	"database/sql"
	"github.com/elliotcourant/noahdb/testutils"
	"github.com/jackc/pgx"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDataGripRegression(t *testing.T) {
	t.Skip("this isnt working yet.")
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
		t.Skip("this is using pgx which isnt yet supported")
		str := testutils.ConnectionString(colony.Addr())
		config, err := pgx.ParseConnectionString(str)
		if !assert.NoError(t, err) {
			panic(err)
		}
		db, err := pgx.Connect(config)
		if err != nil {
			panic(err)
		}
		defer db.Close()

		row := db.QueryRow(`select current_database() as a, current_schemas(false) as b`)
		currentDatabase, currentSchemas := "", ""
		if err := row.Scan(&currentDatabase, &currentSchemas); !assert.NoError(t, err) {
			panic(err)
		}
		assert.NotEmpty(t, currentDatabase)
		assert.NotEmpty(t, currentSchemas)
	})
}
