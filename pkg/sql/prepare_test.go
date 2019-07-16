package sql_test

import (
	"database/sql"
	"fmt"
	"github.com/elliotcourant/noahdb/pkg/core"
	"github.com/elliotcourant/noahdb/pkg/types"
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
				Type:       types.Type_int8,
				PrimaryKey: true,
			},
			{
				ColumnName: "name",
				Type:       types.Type_text,
				PrimaryKey: false,
				Nullable:   true,
			},
		}
		newTable, newColumns, err := colony.Tables().NewTable(table, columns)
		assert.NoError(t, err)
		assert.True(t, newTable.TableID > 0)
		assert.NotEmpty(t, newColumns)
		// prepared, err := db.Prepare("SELECT $1::int, $2::text")
		prepared, err := db.Prepare("SELECT $1::int, $2::int[], a.id, id tenant_id, id::int user_id FROM accounts a WHERE a.id = $3")
		if !assert.NoError(t, err) {
			panic(err)
		}

		input := 1
		row := prepared.QueryRow(input, input+1, input+2)
		fmt.Println(row)
		// value := 0
		// if err := row.Scan(&value); !assert.NoError(t, err) {
		// 	panic(err)
		// }
		// assert.Equal(t, input, value)
	}()
}
