package types_test

import (
	"reflect"
	"testing"

	"github.com/elliotcourant/noahdb/pkg/types"
	"github.com/elliotcourant/noahdb/pkg/types/testutil"
)

func TestCIDTranscode(t *testing.T) {
	typesName := "cid"
	values := []interface{}{
		&types.CID{Uint: 42, Status: types.Present},
		&types.CID{Status: types.Null},
	}
	eqFunc := func(a, b interface{}) bool {
		return reflect.DeepEqual(a, b)
	}

	testutil.TestPgxSuccessfulTranscodeEqFunc(t, typesName, values, eqFunc)

	// No direct conversion from int to cid, convert through text
	testutil.TestPgxSimpleProtocolSuccessfulTranscodeEqFunc(t, "text::"+typesName, values, eqFunc)

	for _, driverName := range []string{"github.com/lib/pq", "github.com/jackc/pgx/stdlib"} {
		testutil.TestDatabaseSQLSuccessfulTranscodeEqFunc(t, driverName, typesName, values, eqFunc)
	}
}

func TestCIDSet(t *testing.T) {
	successfulTests := []struct {
		source interface{}
		result types.CID
	}{
		{source: uint32(1), result: types.CID{Uint: 1, Status: types.Present}},
	}

	for i, tt := range successfulTests {
		var r types.CID
		err := r.Set(tt.source)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if r != tt.result {
			t.Errorf("%d: expected %v to convert to %v, but it was %v", i, tt.source, tt.result, r)
		}
	}
}

func TestCIDAssignTo(t *testing.T) {
	var ui32 uint32
	var pui32 *uint32

	simpleTests := []struct {
		src      types.CID
		dst      interface{}
		expected interface{}
	}{
		{src: types.CID{Uint: 42, Status: types.Present}, dst: &ui32, expected: uint32(42)},
		{src: types.CID{Status: types.Null}, dst: &pui32, expected: (*uint32)(nil)},
	}

	for i, tt := range simpleTests {
		err := tt.src.AssignTo(tt.dst)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if dst := reflect.ValueOf(tt.dst).Elem().Interface(); dst != tt.expected {
			t.Errorf("%d: expected %v to assign %v, but result was %v", i, tt.src, tt.expected, dst)
		}
	}

	pointerAllocTests := []struct {
		src      types.CID
		dst      interface{}
		expected interface{}
	}{
		{src: types.CID{Uint: 42, Status: types.Present}, dst: &pui32, expected: uint32(42)},
	}

	for i, tt := range pointerAllocTests {
		err := tt.src.AssignTo(tt.dst)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if dst := reflect.ValueOf(tt.dst).Elem().Elem().Interface(); dst != tt.expected {
			t.Errorf("%d: expected %v to assign %v, but result was %v", i, tt.src, tt.expected, dst)
		}
	}

	errorTests := []struct {
		src types.CID
		dst interface{}
	}{
		{src: types.CID{Status: types.Null}, dst: &ui32},
	}

	for i, tt := range errorTests {
		err := tt.src.AssignTo(tt.dst)
		if err == nil {
			t.Errorf("%d: expected error but none was returned (%v -> %v)", i, tt.src, tt.dst)
		}
	}
}
