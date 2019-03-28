package types_test

import (
	"github.com/elliotcourant/noahdb/pkg/types/testutil"
	"reflect"
	"testing"

	"github.com/elliotcourant/noahdb/pkg/types"
)

func TestACLItemTranscode(t *testing.T) {
	testutil.TestSuccessfulTranscode(t, "aclitem", []interface{}{
		&types.ACLItem{String: "postgres=arwdDxt/postgres", Status: types.Present},
		&types.ACLItem{String: `postgres=arwdDxt/" tricky, ' } "" \ test user "`, Status: types.Present},
		&types.ACLItem{Status: types.Null},
	})
}

func TestACLItemSet(t *testing.T) {
	successfulTests := []struct {
		source interface{}
		result types.ACLItem
	}{
		{source: "postgres=arwdDxt/postgres", result: types.ACLItem{String: "postgres=arwdDxt/postgres", Status: types.Present}},
		{source: (*string)(nil), result: types.ACLItem{Status: types.Null}},
	}

	for i, tt := range successfulTests {
		var d types.ACLItem
		err := d.Set(tt.source)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if d != tt.result {
			t.Errorf("%d: expected %v to convert to %v, but it was %v", i, tt.source, tt.result, d)
		}
	}
}

func TestACLItemAssignTo(t *testing.T) {
	var s string
	var ps *string

	simpleTests := []struct {
		src      types.ACLItem
		dst      interface{}
		expected interface{}
	}{
		{src: types.ACLItem{String: "postgres=arwdDxt/postgres", Status: types.Present}, dst: &s, expected: "postgres=arwdDxt/postgres"},
		{src: types.ACLItem{Status: types.Null}, dst: &ps, expected: (*string)(nil)},
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
		src      types.ACLItem
		dst      interface{}
		expected interface{}
	}{
		{src: types.ACLItem{String: "postgres=arwdDxt/postgres", Status: types.Present}, dst: &ps, expected: "postgres=arwdDxt/postgres"},
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
		src types.ACLItem
		dst interface{}
	}{
		{src: types.ACLItem{Status: types.Null}, dst: &s},
	}

	for i, tt := range errorTests {
		err := tt.src.AssignTo(tt.dst)
		if err == nil {
			t.Errorf("%d: expected error but none was returned (%v -> %v)", i, tt.src, tt.dst)
		}
	}
}
