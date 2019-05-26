package core_test

import (
	"github.com/elliotcourant/noahdb/pkg/core"
	"github.com/elliotcourant/noahdb/testutils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTypeContext_GetTypeByName(t *testing.T) {
	colony, cleanup := testutils.NewTestColony()
	defer cleanup()

	assertValidType := func(t *testing.T, name string, expected core.Type) {
		typ, ok, err := colony.Types().GetTypeByName(name)
		assert.NoError(t, err)
		assert.True(t, ok)
		assert.Equal(t, expected, typ, "expected %s [%d] found %s [%d]", expected, expected, typ, typ)
	}

	t.Run("int", func(t *testing.T) {
		assertValidType(t, "smallint", core.Type_int2)
		assertValidType(t, "[]smallint", core.Type_int2_array)

		assertValidType(t, "int", core.Type_int4)
		assertValidType(t, "[]int", core.Type_int4_array)

		assertValidType(t, "integer", core.Type_int4)
		assertValidType(t, "[]integer", core.Type_int4_array)

		assertValidType(t, "int8", core.Type_int8)
		assertValidType(t, "[]int8", core.Type_int8_array)

		assertValidType(t, "bigint", core.Type_int8)
		assertValidType(t, "[]bigint", core.Type_int8_array)
	})

	t.Run("text", func(t *testing.T) {
		assertValidType(t, "text", core.Type_text)
		assertValidType(t, "STRING", core.Type_text)
	})

	t.Run("dates and times", func(t *testing.T) {
		assertValidType(t, "timestamp", core.Type_timestamp)
		assertValidType(t, "timestamp without time zone", core.Type_timestamp)
		assertValidType(t, "timestamp with time zone", core.Type_timestamptz)

		assertValidType(t, "timestamp 6", core.Type_timestamp)
		assertValidType(t, "timestamp 5 without time zone", core.Type_timestamp)
		assertValidType(t, "timestamp 4 with time zone", core.Type_timestamptz)

		assertValidType(t, "[]timestamp", core.Type_timestamp_array)
		assertValidType(t, "[]timestamp without time zone", core.Type_timestamp_array)
		assertValidType(t, "[]timestamp with time zone", core.Type_timestamptz_array)

		assertValidType(t, "[]timestamp 6", core.Type_timestamp_array)
		assertValidType(t, "[]timestamp 5 without time zone", core.Type_timestamp_array)
		assertValidType(t, "[]timestamp 4 with time zone", core.Type_timestamptz_array)

		assertValidType(t, "date", core.Type_date)

		assertValidType(t, "[]date", core.Type_date_array)

		assertValidType(t, "time", core.Type_time)
		assertValidType(t, "time without time zone", core.Type_time)
		assertValidType(t, "time with time zone", core.Type_timetz)

		assertValidType(t, "time 6", core.Type_time)
		assertValidType(t, "time 5 without time zone", core.Type_time)
		assertValidType(t, "time 4 with time zone", core.Type_timetz)

		assertValidType(t, "[]time", core.Type_time_array)
		assertValidType(t, "[]time without time zone", core.Type_time_array)
		assertValidType(t, "[]time with time zone", core.Type_timetz_array)

		assertValidType(t, "[]time 6", core.Type_time_array)
		assertValidType(t, "[]time 5 without time zone", core.Type_time_array)
		assertValidType(t, "[]time 4 with time zone", core.Type_timetz_array)
	})

	t.Run("get type array", func(t *testing.T) {
		assertValidType(t, "[]int8", core.Type_int8_array)
	})

	t.Run("get type array with bounds", func(t *testing.T) {
		assertValidType(t, "[12]int8", core.Type_int8_array)
	})
}
