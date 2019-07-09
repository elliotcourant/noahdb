package sql_test

import (
	"database/sql"
	"github.com/elliotcourant/noahdb/testutils"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_ExecPrepare(t *testing.T) {
	t.Skip()
	colony, cleanup := testutils.NewTestColony(t)
	defer cleanup()
	time.Sleep(1 * time.Second)
	func() {
		db, err := sql.Open("postgres", testutils.ConnectionString(colony.Addr()))
		if err != nil {
			panic(err)
		}
		defer db.Close()

		prepared, err := db.Prepare("SELECT $1::int, $2::int[], table_one.id, account_id tenant_id FROM table_one")
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
