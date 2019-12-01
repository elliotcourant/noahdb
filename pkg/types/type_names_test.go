package types

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetTypeByName(t *testing.T) {

	assertValidType := func(t *testing.T, name string, expected Type) {
		typ, ok, err := GetTypeByName(name)
		assert.NoError(t, err)
		assert.True(t, ok, "expect to find %s [%d] for [%s] but the type was not found", expected, expected, name)
		assert.Equal(t, expected, typ, "expected %s [%d] found %s [%d]", expected, expected, typ, typ)
	}

	t.Run("int", func(t *testing.T) {
		assertValidType(t, "smallint", Type_int2)
		assertValidType(t, "smallint[]", Type_int2_array)

		assertValidType(t, "int", Type_int4)
		assertValidType(t, "int[]", Type_int4_array)

		assertValidType(t, "integer", Type_int4)
		assertValidType(t, "integer[]", Type_int4_array)

		assertValidType(t, "int8", Type_int8)
		assertValidType(t, "int8[]", Type_int8_array)

		assertValidType(t, "bigint", Type_int8)
		assertValidType(t, "bigint[]", Type_int8_array)
	})

	t.Run("text", func(t *testing.T) {
		assertValidType(t, "text", Type_text)
		assertValidType(t, "STRING", Type_text)
	})

	t.Run("dates and times", func(t *testing.T) {
		assertValidType(t, "timestamp", Type_timestamp)
		assertValidType(t, "timestamp without time zone", Type_timestamp)
		assertValidType(t, "timestamp with time zone", Type_timestamptz)

		assertValidType(t, "timestamp 6", Type_timestamp)
		assertValidType(t, "timestamp 5 without time zone", Type_timestamp)
		assertValidType(t, "timestamp 4 with time zone", Type_timestamptz)

		assertValidType(t, "timestamp[]", Type_timestamp_array)
		assertValidType(t, "timestamp without time zone[]", Type_timestamp_array)
		assertValidType(t, "timestamp with time zone[]", Type_timestamptz_array)

		assertValidType(t, "timestamp 6[]", Type_timestamp_array)
		assertValidType(t, "timestamp 5 without time zone[]", Type_timestamp_array)
		assertValidType(t, "timestamp 4 with time zone[]", Type_timestamptz_array)

		assertValidType(t, "date", Type_date)

		assertValidType(t, "date[]", Type_date_array)

		assertValidType(t, "time", Type_time)
		assertValidType(t, "time without time zone", Type_time)
		assertValidType(t, "time with time zone", Type_timetz)

		assertValidType(t, "time 6", Type_time)
		assertValidType(t, "time 5 without time zone", Type_time)
		assertValidType(t, "time 4 with time zone", Type_timetz)

		assertValidType(t, "time[]", Type_time_array)
		assertValidType(t, "time without time zone[]", Type_time_array)
		assertValidType(t, "time with time zone[]", Type_timetz_array)

		assertValidType(t, "time 6[]", Type_time_array)
		assertValidType(t, "time 5 without time zone[]", Type_time_array)
		assertValidType(t, "time 4 with time zone[]", Type_timetz_array)
	})

	t.Run("get type array", func(t *testing.T) {
		assertValidType(t, "int8[]", Type_int8_array)
	})

	t.Run("get type array with bounds", func(t *testing.T) {
		assertValidType(t, "int8[12]", Type_int8_array)
	})
}

func TestTypeContext_GetTypeByOid(t *testing.T) {
	assertCorrectType := func(oid OID, expected Type) {
		result, ok := GetTypeByOid(oid)
		if !assert.True(t, ok, "could not find matching type") {
			t.FailNow()
		}
		if !assert.Equal(t, expected, result) {
			t.FailNow()
		}
	}

	assertMissingType := func(oid OID) {
		_, ok := GetTypeByOid(oid)
		if !assert.False(t, ok, "found matching type") {
			t.FailNow()
		}
	}

	t.Run("general oid to type", func(t *testing.T) {
		assertCorrectType(BoolOID, Type_bool)
		assertCorrectType(ByteaOID, Type_bytea)
		assertCorrectType(CharOID, Type_char)
		assertCorrectType(Int8OID, Type_int8)
	})

	t.Run("missing types", func(t *testing.T) {
		assertMissingType(12512121)
		assertMissingType(1)
	})
}
