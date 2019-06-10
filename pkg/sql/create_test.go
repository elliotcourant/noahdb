package sql_test

import (
	"database/sql"
	"github.com/elliotcourant/noahdb/testutils"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewCreateStatementPlan(t *testing.T) {
	colony, cleanup := testutils.NewPgTestColony(t)
	defer cleanup()

	db, err := sql.Open("postgres", testutils.ConnectionString(colony.Addr()))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	t.Run("create simple table", func(t *testing.T) {
		_, err = db.Exec(`CREATE TABLE accounts (id BIGSERIAL PRIMARY KEY, name TEXT UNIQUE);`)
		assert.NoError(t, err)
		table, ok, err := colony.Tables().GetTable("accounts")
		assert.NoError(t, err)
		assert.True(t, ok)
		assert.NotEmpty(t, table)
	})

	t.Run("create tenants table", func(t *testing.T) {
		_, err = db.Exec(`CREATE TABLE accounts_tenants (id BIGSERIAL PRIMARY KEY, name TEXT NULL UNIQUE) TABLESPACE "noah.tenants";`)
		assert.NoError(t, err)
	})
}
