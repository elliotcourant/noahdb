package pgwire_test

import (
	"database/sql"
	"github.com/elliotcourant/noahdb/testutils"
	_ "github.com/lib/pq"
	"github.com/readystock/golog"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestLibPqStartup(t *testing.T) {
	colony, cleanup := testutils.NewTestColony(t)
	defer cleanup()
	time.Sleep(5 * time.Second)
	func() {
		db, err := sql.Open("postgres", testutils.ConnectionString(colony.Addr()))
		if err != nil {
			panic(err)
		}
		defer db.Close()

		if row, err := db.Query(`SELECT 1;`); err != nil {
			panic(err)
		} else {
			row.Next()
			intVal := 0
			row.Scan(&intVal)
			assert.Equal(t, 1, intVal)
			row.Close()
		}

		start := time.Now()
		if row, err := db.Query(`SELECT 1;`); err != nil {
			panic(err)
		} else {
			row.Close()
		}
		golog.Infof("query time: %s", time.Since(start))
	}()
}

func Test_HandleParse_BadSyntax(t *testing.T) {
	colony, cleanup := testutils.NewTestColony(t)
	defer cleanup()
	time.Sleep(1 * time.Second)
	func() {
		db, err := sql.Open("postgres", testutils.ConnectionString(colony.Addr()))
		if err != nil {
			panic(err)
		}
		defer db.Close()
		_, err = db.Query(`SELECTad;`)
		assert.Error(t, err) // There should not be an error sending the query.
	}()
}
