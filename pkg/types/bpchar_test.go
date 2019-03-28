package types_test

import (
	"reflect"
	"testing"

	"github.com/elliotcourant/noahdb/pkg/types"
	"github.com/elliotcourant/noahdb/pkg/types/testutil"
)

func TestChar3Transcode(t *testing.T) {
	testutil.TestSuccessfulTranscodeEqFunc(t, "char(3)", []interface{}{
		&types.BPChar{String: "a  ", Status: types.Present},
		&types.BPChar{String: " a ", Status: types.Present},
		&types.BPChar{String: "嗨  ", Status: types.Present},
		&types.BPChar{String: "   ", Status: types.Present},
		&types.BPChar{Status: types.Null},
	}, func(aa, bb interface{}) bool {
		a := aa.(types.BPChar)
		b := bb.(types.BPChar)

		return a.Status == b.Status && a.String == b.String
	})
}

func TestBPCharAssignTo(t *testing.T) {
	var (
		str string
		run rune
	)
	simpleTests := []struct {
		src      types.BPChar
		dst      interface{}
		expected interface{}
	}{
		{src: types.BPChar{String: "simple", Status: types.Present}, dst: &str, expected: "simple"},
		{src: types.BPChar{String: "嗨", Status: types.Present}, dst: &run, expected: '嗨'},
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

}
