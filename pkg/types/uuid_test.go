package types_test

import (
	"bytes"
	"testing"

	"github.com/elliotcourant/noahdb/pkg/types"
	"github.com/elliotcourant/noahdb/pkg/types/testutil"
)

func TestUUIDTranscode(t *testing.T) {
	testutil.TestSuccessfulTranscode(t, "uuid", []interface{}{
		&types.UUID{Bytes: [16]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, Status: types.Present},
		&types.UUID{Status: types.Null},
	})
}

func TestUUIDSet(t *testing.T) {
	successfulTests := []struct {
		source interface{}
		result types.UUID
	}{
		{
			source: nil,
			result: types.UUID{Status: types.Null},
		},
		{
			source: [16]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
			result: types.UUID{Bytes: [16]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, Status: types.Present},
		},
		{
			source: []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
			result: types.UUID{Bytes: [16]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, Status: types.Present},
		},
		{
			source: ([]byte)(nil),
			result: types.UUID{Status: types.Null},
		},
		{
			source: "00010203-0405-0607-0809-0a0b0c0d0e0f",
			result: types.UUID{Bytes: [16]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, Status: types.Present},
		},
	}

	for i, tt := range successfulTests {
		var r types.UUID
		err := r.Set(tt.source)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if r != tt.result {
			t.Errorf("%d: expected %v to convert to %v, but it was %v", i, tt.source, tt.result, r)
		}
	}
}

func TestUUIDAssignTo(t *testing.T) {
	{
		src := types.UUID{Bytes: [16]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, Status: types.Present}
		var dst [16]byte
		expected := [16]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}

		err := src.AssignTo(&dst)
		if err != nil {
			t.Error(err)
		}

		if dst != expected {
			t.Errorf("expected %v to assign %v, but result was %v", src, expected, dst)
		}
	}

	{
		src := types.UUID{Bytes: [16]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, Status: types.Present}
		var dst []byte
		expected := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}

		err := src.AssignTo(&dst)
		if err != nil {
			t.Error(err)
		}

		if bytes.Compare(dst, expected) != 0 {
			t.Errorf("expected %v to assign %v, but result was %v", src, expected, dst)
		}
	}

	{
		src := types.UUID{Bytes: [16]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, Status: types.Present}
		var dst string
		expected := "00010203-0405-0607-0809-0a0b0c0d0e0f"

		err := src.AssignTo(&dst)
		if err != nil {
			t.Error(err)
		}

		if dst != expected {
			t.Errorf("expected %v to assign %v, but result was %v", src, expected, dst)
		}
	}

}
