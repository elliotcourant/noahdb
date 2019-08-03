package core_test

import (
	"database/sql"
	"github.com/elliotcourant/noahdb/testutils"
	"github.com/stretchr/testify/assert"
	"testing"
)

type (
	regproc = string
	oid     = int
)

func TestPostgres(t *testing.T) {
	colony, cleanup := testutils.NewPgTestColony(t)
	defer cleanup()
	t.Run("aggregates", func(t *testing.T) {
		db, err := sql.Open("postgres", testutils.ConnectionString(colony.Addr()))
		if err != nil {
			panic(err)
		}
		defer db.Close()

		rows, err := db.Query(`SELECT * FROM pg_catalog.pg_aggregate;`)
		assert.NoError(t, err)

		for rows.Next() {
			item := struct {
				aggfnoid         string
				aggkind          string
				aggnumdirectargs int
				aggtransfn       regproc
				aggfinalfn       regproc
				aggcombinefn     regproc
				aggserialfn      regproc
				aggdeserialfn    regproc
				aggmtransfn      regproc
				aggminvtransfn   regproc
				aggmfinalfn      regproc
				aggfinalextra    bool
				aggmfinalextra   bool
				aggfinalmodify   string
				aggmfinalmodify  string
				aggsortop        oid
				aggtranstype     oid
				aggtransspace    int
				aggmtranstype    oid
				aggmtransspace   int
				agginitval       *string
				aggminitval      *string
			}{}
			err := rows.Scan(
				&item.aggfnoid,
				&item.aggkind,
				&item.aggnumdirectargs,
				&item.aggtransfn,
				&item.aggfinalfn,
				&item.aggcombinefn,
				&item.aggserialfn,
				&item.aggdeserialfn,
				&item.aggmtransfn,
				&item.aggminvtransfn,
				&item.aggmfinalfn,
				&item.aggfinalextra,
				&item.aggmfinalextra,
				&item.aggfinalmodify,
				&item.aggmfinalmodify,
				&item.aggsortop,
				&item.aggtranstype,
				&item.aggtransspace,
				&item.aggmtranstype,
				&item.aggmtransspace,
				&item.agginitval,
				&item.aggminitval,
			)
			assert.NoError(t, err)
			assert.NotEmpty(t, item)
		}

		err = rows.Err()
		assert.NoError(t, err)
	})

	t.Run("access methods", func(t *testing.T) {
		db, err := sql.Open("postgres", testutils.ConnectionString(colony.Addr()))
		if err != nil {
			panic(err)
		}
		defer db.Close()

		rows, err := db.Query(`SELECT * FROM pg_catalog.pg_am;`)
		assert.NoError(t, err)

		for rows.Next() {
			item := struct {
				amname    string
				amhandler regproc
				amtype    string
			}{}
			err := rows.Scan(
				&item.amname,
				&item.amhandler,
				&item.amtype,
			)
			assert.NoError(t, err)
			assert.NotEmpty(t, item)
		}

		err = rows.Err()
		assert.NoError(t, err)
	})
}
