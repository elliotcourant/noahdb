package types_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/elliotcourant/noahdb/pkg/types"
	"github.com/elliotcourant/noahdb/pkg/types/testutil"
	"github.com/readystock/pgx"
)

func TestRecordTranscode(t *testing.T) {
	t.Skip("cannot transcode from pgx into npgx")
	conn := testutil.MustConnectPgx(t)
	defer testutil.MustClose(t, conn)

	tests := []struct {
		sql      string
		expected types.Record
	}{
		{
			sql: `select row()`,
			expected: types.Record{
				Fields: []types.Value{},
				Status: types.Present,
			},
		},
		{
			sql: `select row('foo'::text, 42::int4)`,
			expected: types.Record{
				Fields: []types.Value{
					&types.Text{String: "foo", Status: types.Present},
					&types.Int4{Int: 42, Status: types.Present},
				},
				Status: types.Present,
			},
		},
		{
			sql: `select row(100.0::float4, 1.09::float4)`,
			expected: types.Record{
				Fields: []types.Value{
					&types.Float4{Float: 100, Status: types.Present},
					&types.Float4{Float: 1.09, Status: types.Present},
				},
				Status: types.Present,
			},
		},
		{
			sql: `select row('foo'::text, array[1, 2, null, 4]::int4[], 42::int4)`,
			expected: types.Record{
				Fields: []types.Value{
					&types.Text{String: "foo", Status: types.Present},
					&types.Int4Array{
						Elements: []types.Int4{
							{Int: 1, Status: types.Present},
							{Int: 2, Status: types.Present},
							{Status: types.Null},
							{Int: 4, Status: types.Present},
						},
						Dimensions: []types.ArrayDimension{{Length: 4, LowerBound: 1}},
						Status:     types.Present,
					},
					&types.Int4{Int: 42, Status: types.Present},
				},
				Status: types.Present,
			},
		},
		{
			sql: `select row(null)`,
			expected: types.Record{
				Fields: []types.Value{
					&types.Unknown{Status: types.Null},
				},
				Status: types.Present,
			},
		},
		{
			sql: `select null::record`,
			expected: types.Record{
				Status: types.Null,
			},
		},
	}

	for i, tt := range tests {
		psName := fmt.Sprintf("test%d", i)
		ps, err := conn.Prepare(psName, tt.sql)
		if err != nil {
			t.Fatal(err)
		}
		ps.FieldDescriptions[0].FormatCode = pgx.BinaryFormatCode

		var result types.Record
		if err := conn.QueryRow(psName).Scan(&result); err != nil {
			t.Errorf("%d: %v", i, err)
			continue
		}

		if !reflect.DeepEqual(tt.expected, result) {
			t.Errorf("%d: expected %#v, got %#v", i, tt.expected, result)
		}
	}
}

func TestRecordWithUnknownOID(t *testing.T) {
	conn := testutil.MustConnectPgx(t)
	defer testutil.MustClose(t, conn)

	_, err := conn.Exec(`drop type if exists floatrange;

create type floatrange as range (
  subtype = float8,
  subtype_diff = float8mi
);`)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Exec("drop type floatrange")

	var result types.Record
	err = conn.QueryRow("select row('foo'::text, floatrange(1, 10), 'bar'::text)").Scan(&result)
	if err == nil {
		t.Errorf("expected error but none")
	}
}

func TestRecordAssignTo(t *testing.T) {
	var valueSlice []types.Value
	var interfaceSlice []interface{}

	simpleTests := []struct {
		src      types.Record
		dst      interface{}
		expected interface{}
	}{
		{
			src: types.Record{
				Fields: []types.Value{
					&types.Text{String: "foo", Status: types.Present},
					&types.Int4{Int: 42, Status: types.Present},
				},
				Status: types.Present,
			},
			dst: &valueSlice,
			expected: []types.Value{
				&types.Text{String: "foo", Status: types.Present},
				&types.Int4{Int: 42, Status: types.Present},
			},
		},
		{
			src: types.Record{
				Fields: []types.Value{
					&types.Text{String: "foo", Status: types.Present},
					&types.Int4{Int: 42, Status: types.Present},
				},
				Status: types.Present,
			},
			dst:      &interfaceSlice,
			expected: []interface{}{"foo", int32(42)},
		},
		{
			src:      types.Record{Status: types.Null},
			dst:      &valueSlice,
			expected: ([]types.Value)(nil),
		},
		{
			src:      types.Record{Status: types.Null},
			dst:      &interfaceSlice,
			expected: ([]interface{})(nil),
		},
	}

	for i, tt := range simpleTests {
		err := tt.src.AssignTo(tt.dst)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if dst := reflect.ValueOf(tt.dst).Elem().Interface(); !reflect.DeepEqual(dst, tt.expected) {
			t.Errorf("%d: expected %v to assign %v, but result was %v", i, tt.src, tt.expected, dst)
		}
	}
}
