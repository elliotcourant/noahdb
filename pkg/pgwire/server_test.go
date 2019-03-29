package pgwire_test

import (
	"database/sql"
	"fmt"
	"github.com/elliotcourant/noahdb/pkg/pgwire"
	_ "github.com/lib/pq"
	"github.com/readystock/golog"
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
	"time"
)

type config struct {
	address string
	port    int
}

func (conf config) Address() string {
	return conf.address
}

func (conf config) Port() int {
	return conf.port
}

func (conf config) LibPqConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		conf.Address(), conf.Port(), "noah", "password", "postgres")
}

func NewConfig() config {
	ln, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		panic(err)
	}

	defer ln.Close()
	addr, err := net.ResolveTCPAddr(ln.Addr().Network(), ln.Addr().String())
	if err != nil {
		panic(err)
	}

	return config{
		address: addr.IP.String(),
		port:    addr.Port,
	}
}

func TestLibPqStartup(t *testing.T) {
	conf := NewConfig()
	go func() {
		if err := pgwire.NewServer(conf); err != nil {
			panic(err)
		}
	}()
	time.Sleep(1 * time.Second)
	for i := 0; i < 5; i++ {
		func() {
			db, err := sql.Open("postgres", conf.LibPqConnectionString())
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
}
