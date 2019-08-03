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
	t.Run("aggregate", func(t *testing.T) {
		db, err := sql.Open("postgres", testutils.ConnectionString(colony.Addr()))
		if err != nil {
			panic(err)
		}
		defer db.Close()

		rows, err := db.Query(`SELECT * FROM pg_catalog.pg_aggregate;`)
		assert.NoError(t, err)

		for rows.Next() {
			aggregate := struct {
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
				&aggregate.aggfnoid,
				&aggregate.aggkind,
				&aggregate.aggnumdirectargs,
				&aggregate.aggtransfn,
				&aggregate.aggfinalfn,
				&aggregate.aggcombinefn,
				&aggregate.aggserialfn,
				&aggregate.aggdeserialfn,
				&aggregate.aggmtransfn,
				&aggregate.aggminvtransfn,
				&aggregate.aggmfinalfn,
				&aggregate.aggfinalextra,
				&aggregate.aggmfinalextra,
				&aggregate.aggfinalmodify,
				&aggregate.aggmfinalmodify,
				&aggregate.aggsortop,
				&aggregate.aggtranstype,
				&aggregate.aggtransspace,
				&aggregate.aggmtranstype,
				&aggregate.aggmtransspace,
				&aggregate.agginitval,
				&aggregate.aggminitval,
			)
			assert.NoError(t, err)
			assert.NotEmpty(t, aggregate)
		}

		err = rows.Err()
		assert.NoError(t, err)
	})
}
