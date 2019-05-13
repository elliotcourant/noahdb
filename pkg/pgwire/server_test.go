package pgwire_test

import (
	"database/sql"
	"fmt"
	"github.com/elliotcourant/noahdb/pkg/core"
	"github.com/elliotcourant/noahdb/pkg/pgwire"
	_ "github.com/lib/pq"
	"github.com/readystock/golog"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net"
	"os"
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

func NewColony() (core.Colony, core.TransportWrapper, func()) {
	tempdir, err := ioutil.TempDir("", "core-temp")
	if err != nil {
		panic(err)
	}

	colony, wrapper, err := core.NewColony(tempdir, "", ":")
	if err != nil {
		panic(err)
	}

	return colony, wrapper, func() {
		if err := os.RemoveAll(tempdir); err != nil {
			panic(err)
		}
	}
}

func TestLibPqStartup(t *testing.T) {
	colony, wrapper, cleanup := NewColony()
	defer cleanup()
	go func() {
		if err := pgwire.NewServer(colony, wrapper); err != nil {
			panic(err)
		}
	}()
	time.Sleep(1 * time.Second)
	func() {
		db, err := sql.Open("postgres", LibPqConnectionString(wrapper.Addr()))
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
