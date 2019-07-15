package core_test

import (
	"github.com/elliotcourant/noahdb/pkg/types"
	"github.com/elliotcourant/noahdb/testutils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTypeContext_GetTypeByName(t *testing.T) {
	colony, cleanup := testutils.NewTestColony(t)
	defer cleanup()

	assertValidType := func(t *testing.T, name string, expected types.Type) {
		typ, ok, err := colony.Types().GetTypeByName(name)
		assert.NoError(t, err)
		assert.True(t, ok)
		assert.Equal(t, expected, typ, "expected %s [%d] found %s [%d]", expected, expected, typ, typ)
	}

	t.Run("int", func(t *testing.T) {
		assertValidType(t, "smallint", types.Type_int2)
		assertValidType(t, "smallint[]", types.Type_int2_array)

		assertValidType(t, "int", types.Type_int4)
		assertValidType(t, "int[]", types.Type_int4_array)

		assertValidType(t, "integer", types.Type_int4)
		assertValidType(t, "integer[]", types.Type_int4_array)

		assertValidType(t, "int8", types.Type_int8)
		assertValidType(t, "int8[]", types.Type_int8_array)

		assertValidType(t, "bigint", types.Type_int8)
		assertValidType(t, "bigint[]", types.Type_int8_array)
	})

	t.Run("text", func(t *testing.T) {
		assertValidType(t, "text", types.Type_text)
		assertValidType(t, "STRING", types.Type_text)
	})

	t.Run("dates and times", func(t *testing.T) {
		assertValidType(t, "timestamp", types.Type_timestamp)
		assertValidType(t, "timestamp without time zone", types.Type_timestamp)
		assertValidType(t, "timestamp with time zone", types.Type_timestamptz)

		assertValidType(t, "timestamp 6", types.Type_timestamp)
		assertValidType(t, "timestamp 5 without time zone", types.Type_timestamp)
		assertValidType(t, "timestamp 4 with time zone", types.Type_timestamptz)

		assertValidType(t, "timestamp[]", types.Type_timestamp_array)
		assertValidType(t, "timestamp without time zone[]", types.Type_timestamp_array)
		assertValidType(t, "timestamp with time zone[]", types.Type_timestamptz_array)

		assertValidType(t, "timestamp 6[]", types.Type_timestamp_array)
		assertValidType(t, "timestamp 5 without time zone[]", types.Type_timestamp_array)
		assertValidType(t, "timestamp 4 with time zone[]", types.Type_timestamptz_array)

		assertValidType(t, "date", types.Type_date)

		assertValidType(t, "date[]", types.Type_date_array)

		assertValidType(t, "time", types.Type_time)
		assertValidType(t, "time without time zone", types.Type_time)
		assertValidType(t, "time with time zone", types.Type_timetz)

		assertValidType(t, "time 6", types.Type_time)
		assertValidType(t, "time 5 without time zone", types.Type_time)
		assertValidType(t, "time 4 with time zone", types.Type_timetz)

		assertValidType(t, "time[]", types.Type_time_array)
		assertValidType(t, "time without time zone[]", types.Type_time_array)
		assertValidType(t, "time with time zone[]", types.Type_timetz_array)

		assertValidType(t, "time 6[]", types.Type_time_array)
		assertValidType(t, "time 5 without time zone[]", types.Type_time_array)
		assertValidType(t, "time 4 with time zone[]", types.Type_timetz_array)
	})

	t.Run("get type array", func(t *testing.T) {
		assertValidType(t, "int8[]", types.Type_int8_array)
	})

	t.Run("get type array with bounds", func(t *testing.T) {
		assertValidType(t, "int8[12]", types.Type_int8_array)
	})
}

func TestTypeContext_GetTypeByOid(t *testing.T) {
	colony, cleanup := testutils.NewTestColony(t)
	defer cleanup()

	assertCorrectType := func(oid types.OID, expected types.Type) {
		result, ok := colony.Types().GetTypeByOid(oid)
		if !assert.True(t, ok, "could not find matching type") {
			t.FailNow()
		}
		if !assert.Equal(t, expected, result) {
			t.FailNow()
		}
	}

	assertMissingType := func(oid types.OID) {
		_, ok := colony.Types().GetTypeByOid(oid)
		if !assert.False(t, ok, "found matching type") {
			t.FailNow()
		}
	}

	t.Run("general oid to type", func(t *testing.T) {
		assertCorrectType(types.BoolOID, types.Type_bool)
		assertCorrectType(types.ByteaOID, types.Type_bytea)
		assertCorrectType(types.CharOID, types.Type_char)
		assertCorrectType(types.Int8OID, types.Type_int8)
	})

	t.Run("missing types", func(t *testing.T) {
		assertMissingType(12512121)
		assertMissingType(1)
	})
}
