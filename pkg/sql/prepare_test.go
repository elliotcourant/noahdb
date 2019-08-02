package sql_test

import (
	"database/sql"
	"fmt"
	"github.com/elliotcourant/noahdb/testutils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ExecPrepare(t *testing.T) {
	colony, cleanup := testutils.NewPgTestColony(t)
	defer cleanup()
	func() {
		// db, err := sql.Open("postgres", "postgresql://postgres@localhost:5432?sslmode=disable")
		db, err := sql.Open("postgres", testutils.ConnectionString(colony.Addr()))
		if err != nil {
			panic(err)
		}
		defer db.Close()

		prepared, err := db.Prepare("SELECT $1::int, $2::int")
		if !assert.NoError(t, err) {
			panic(err)
		}

		input := 1
		row := prepared.QueryRow(input, input+1)
		arg1, arg2 := 0, 0
		if err := row.Scan(&arg1, &arg2); err != nil {
			panic(err)
		}
		fmt.Println(arg1, arg2)
	}()
}
