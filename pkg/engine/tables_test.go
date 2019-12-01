package engine_test

import (
	"github.com/elliotcourant/noahdb/pkg/engine"
	"github.com/elliotcourant/noahdb/pkg/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTableContextBase_NewTable(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		cluster, cleanup := NewTestCoreCluster(t, 1)
		defer cleanup()

		txn := cluster.Begin(t)

		table, columns, err := txn.Tables().NewTable(
			engine.Table{
				Schema:      "public",
				Name:        "table",
				Type:        engine.TableType_Master,
				HasSequence: false,
			},
			[]engine.Column{
				{
					Type:            types.Type_int8,
					Index:           0,
					Name:            "table_id",
					IsPrimaryKey:    true,
					IsNullable:      false,
					IsShardKey:      false,
					IsSerial:        false,
					ForeignColumnId: 0,
				},
				{
					Type:            types.Type_text,
					Index:           1,
					Name:            "name",
					IsPrimaryKey:    false,
					IsNullable:      false,
					IsShardKey:      false,
					IsSerial:        false,
					ForeignColumnId: 0,
				},
			})
		assert.NoError(t, err)
		assert.NotZero(t, table.TableId)
		assert.NotEmpty(t, columns)
	})
}
