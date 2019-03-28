package types_test

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/elliotcourant/noahdb/pkg/types"
	"github.com/elliotcourant/noahdb/pkg/types/testutil"
)

func TestTextTranscode(t *testing.T) {
	for _, typesName := range []string{"text", "varchar"} {
		testutil.TestSuccessfulTranscode(t, typesName, []interface{}{
			&types.Text{String: "", Status: types.Present},
			&types.Text{String: "foo", Status: types.Present},
			&types.Text{Status: types.Null},
		})
	}
}

func TestTextSet(t *testing.T) {
	successfulTests := []struct {
		source interface{}
		result types.Text
	}{
		{source: "foo", result: types.Text{String: "foo", Status: types.Present}},
		{source: _string("bar"), result: types.Text{String: "bar", Status: types.Present}},
		{source: (*string)(nil), result: types.Text{Status: types.Null}},
	}

	for i, tt := range successfulTests {
		var d types.Text
		err := d.Set(tt.source)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if d != tt.result {
			t.Errorf("%d: expected %v to convert to %v, but it was %v", i, tt.source, tt.result, d)
		}
	}
}

func TestTextAssignTo(t *testing.T) {
	var s string
	var ps *string

	stringTests := []struct {
		src      types.Text
		dst      interface{}
		expected interface{}
	}{
		{src: types.Text{String: "foo", Status: types.Present}, dst: &s, expected: "foo"},
		{src: types.Text{Status: types.Null}, dst: &ps, expected: (*string)(nil)},
	}

	for i, tt := range stringTests {
		err := tt.src.AssignTo(tt.dst)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if dst := reflect.ValueOf(tt.dst).Elem().Interface(); dst != tt.expected {
			t.Errorf("%d: expected %v to assign %v, but result was %v", i, tt.src, tt.expected, dst)
		}
	}

	var buf []byte

	bytesTests := []struct {
		src      types.Text
		dst      *[]byte
		expected []byte
	}{
		{src: types.Text{String: "foo", Status: types.Present}, dst: &buf, expected: []byte("foo")},
		{src: types.Text{Status: types.Null}, dst: &buf, expected: nil},
	}

	for i, tt := range bytesTests {
		err := tt.src.AssignTo(tt.dst)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if bytes.Compare(*tt.dst, tt.expected) != 0 {
			t.Errorf("%d: expected %v to assign %v, but result was %v", i, tt.src, tt.expected, tt.dst)
		}
	}

	pointerAllocTests := []struct {
		src      types.Text
		dst      interface{}
		expected interface{}
	}{
		{src: types.Text{String: "foo", Status: types.Present}, dst: &ps, expected: "foo"},
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
		src types.Text
		dst interface{}
	}{
		{src: types.Text{Status: types.Null}, dst: &s},
	}

	for i, tt := range errorTests {
		err := tt.src.AssignTo(tt.dst)
		if err == nil {
			t.Errorf("%d: expected error but none was returned (%v -> %v)", i, tt.src, tt.dst)
		}
	}
}
