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

	db, err := sql.Open("postgres", testutils.ConnectionString(colony.Addr()))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	tryCatalogQuery := func(t *testing.T, query string, scan func(rows *sql.Rows) error) {
		rows, err := db.Query(query)
		assert.NoError(t, err)

		for rows.Next() {
			err := scan(rows)
			assert.NoError(t, err)
		}

		err = rows.Err()
		assert.NoError(t, err)
	}

	t.Run("aggregates", func(t *testing.T) {
		tryCatalogQuery(t, `SELECT * FROM pg_catalog.pg_aggregate;`, func(rows *sql.Rows) error {
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
			return err
		})
	})

	t.Run("access methods", func(t *testing.T) {
		tryCatalogQuery(t, `SELECT * FROM pg_catalog.pg_am;`, func(rows *sql.Rows) error {
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
			return err
		})
	})

	t.Run("access method operators", func(t *testing.T) {
		tryCatalogQuery(t, `SELECT * FROM pg_catalog.pg_amop;`, func(rows *sql.Rows) error {
			item := struct {
				amopfamily     oid
				ampolefttype   oid
				amoprighttype  oid
				amopstrategy   int
				amoppurpose    string
				amopopr        oid
				amopmethod     oid
				amopsortfamily oid
			}{}
			err := rows.Scan(
				&item.amopfamily,
				&item.ampolefttype,
				&item.amoprighttype,
				&item.amopstrategy,
				&item.amoppurpose,
				&item.amopopr,
				&item.amopmethod,
				&item.amopsortfamily,
			)
			assert.NoError(t, err)
			assert.NotEmpty(t, item)
			return err
		})
	})

	t.Run("access method operator support functions", func(t *testing.T) {
		tryCatalogQuery(t, `SELECT * FROM pg_catalog.pg_amproc;`, func(rows *sql.Rows) error {
			item := struct {
				amprocfamily    oid
				amproclefttype  oid
				amprocrighttype oid
				amprocnum       int
				amproc          regproc
			}{}
			err := rows.Scan(
				&item.amprocfamily,
				&item.amproclefttype,
				&item.amprocrighttype,
				&item.amprocnum,
				&item.amproc,
			)
			assert.NoError(t, err)
			assert.NotEmpty(t, item)
			return err
		})
	})
}
