package types_test

import (
	"testing"

	"github.com/elliotcourant/noahdb/pkg/types"
	"github.com/elliotcourant/noahdb/pkg/types/testutil"
)

func TestPointTranscode(t *testing.T) {
	testutil.TestSuccessfulTranscode(t, "point", []interface{}{
		&types.Point{P: types.Vec2{1.234, 5.6789012345}, Status: types.Present},
		&types.Point{P: types.Vec2{-1.234, -5.6789}, Status: types.Present},
		&types.Point{Status: types.Null},
	})
}
