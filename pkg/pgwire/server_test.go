package pgwire_test

import (
	"database/sql"
	"fmt"
	"github.com/elliotcourant/noahdb/testutils"
	_ "github.com/lib/pq"
	"github.com/readystock/golog"
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
	"time"
)

func LibPqConnectionString(address net.Addr) string {
	addr, err := net.ResolveTCPAddr(address.Network(), address.String())
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		addr.IP.String(), addr.Port, "noah", "password", "postgres")
}

func TestLibPqStartup(t *testing.T) {
	t.Skip("shits broke atm")
	colony, cleanup := testutils.NewTestColony()
	defer cleanup()
	time.Sleep(1 * time.Second)
	func() {
		db, err := sql.Open("postgres", LibPqConnectionString(colony.Addr()))
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
	colony, cleanup := testutils.NewTestColony()
	defer cleanup()
	time.Sleep(1 * time.Second)
	func() {
		db, err := sql.Open("postgres", LibPqConnectionString(colony.Addr()))
		if err != nil {
			panic(err)
		}
		defer db.Close()
		_, err = db.Query(`SELECTad;`)
		assert.Error(t, err) // There should not be an error sending the query.
	}()
}
