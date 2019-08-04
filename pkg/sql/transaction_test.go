package sql_test

import (
	"database/sql"
	"github.com/elliotcourant/noahdb/testutils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Transaction(t *testing.T) {
	t.Skip("not finished yet")
	colony, cleanup := testutils.NewPgTestColony(t)
	defer cleanup()

	t.Run("begin and commit", func(t *testing.T) {
		db, err := sql.Open("postgres", testutils.ConnectionString(colony.Addr()))
		if err != nil {
			panic(err)
		}
		defer db.Close()

		tx, err := db.Begin()
		assert.NoError(t, err)

		_, err = tx.Exec(`CREATE TABLE accounts (id BIGSERIAL PRIMARY KEY, name TEXT) TABLESPACE "noah.tenants"`)
		if !assert.NoError(t, err) {
			panic(err)
		}

		err = tx.Commit()
		assert.NoError(t, err)
	})
}
