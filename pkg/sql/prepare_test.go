package sql_test

import (
	"database/sql"
	"github.com/elliotcourant/noahdb/pkg/core"
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
		// db, err := sql.Open("postgres", "postgresql://postgres@localhost:5432?sslmode=disable")
		db, err := sql.Open("postgres", testutils.ConnectionString(colony.Addr()))
		if err != nil {
			panic(err)
		}
		defer db.Close()

		table := core.Table{
			TableName: "accounts",
			TableType: core.TableType_Tenant,
		}
		columns := []core.Column{
			{
				ColumnName: "id",
				Type:       core.Type_int8,
				PrimaryKey: true,
			},
			{
				ColumnName: "name",
				Type:       core.Type_text,
				PrimaryKey: false,
				Nullable:   true,
			},
		}
		newTable, newColumns, err := colony.Tables().NewTable(table, columns)
		assert.NoError(t, err)
		assert.True(t, newTable.TableID > 0)
		assert.NotEmpty(t, newColumns)
		prepared, err := db.Prepare("SELECT $1::int, $2::bigint")
		// prepared, err := db.Prepare("SELECT $1::int, $2::int[], a.id, id tenant_id, id::int user_id FROM accounts a")
		if !assert.NoError(t, err) {
			panic(err)
		}

		input := 1
		row := prepared.QueryRow(input, input)

		value := 0
		if err := row.Scan(&value); !assert.NoError(t, err) {
			panic(err)
		}
		assert.Equal(t, input, value)
	}()
}
