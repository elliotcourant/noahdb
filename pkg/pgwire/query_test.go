package pgwire_test

import (
	"database/sql"
	"github.com/elliotcourant/noahdb/testutils"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_Local_Prepared(t *testing.T) {
	t.Skip()
	func() {
		db, err := sql.Open("postgres", "postgresql://postgres@localhost:5432?sslmode=disable")
		if err != nil {
			panic(err)
		}
		defer db.Close()

		prepared, err := db.Prepare("SELECT * FROM account_users WHERE account_id = $1")
		if !assert.NoError(t, err) {
			panic(err)
		}

		input := 1
		row := prepared.QueryRow(input)

		value := 0
		if err := row.Scan(&value); !assert.NoError(t, err) {
			panic(err)
		}
		assert.Equal(t, input, value)
	}()
}

func Test_HandleParse_Prepared(t *testing.T) {
	colony, cleanup := testutils.NewPgTestColony(t)
	defer cleanup()
	time.Sleep(1 * time.Second)
	func() {
		db, err := sql.Open("postgres", testutils.ConnectionString(colony.Addr()))
		if err != nil {
			panic(err)
		}
		defer db.Close()

		prepared, err := db.Prepare("SELECT $1")
		if !assert.NoError(t, err) {
			panic(err)
		}

		input := 1
		row := prepared.QueryRow(input)

		value := 0
		if err := row.Scan(&value); !assert.NoError(t, err) {
			panic(err)
		}
		assert.Equal(t, input, value)
	}()
}
