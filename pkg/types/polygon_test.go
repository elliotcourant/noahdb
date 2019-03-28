package types_test

import (
	"testing"

	"github.com/elliotcourant/noahdb/pkg/types"
	"github.com/elliotcourant/noahdb/pkg/types/testutil"
)

func TestPolygonTranscode(t *testing.T) {
	testutil.TestSuccessfulTranscode(t, "polygon", []interface{}{
		&types.Polygon{
			P:      []types.Vec2{{3.14, 1.678901234}, {7.1, 5.234}, {5.0, 3.234}},
			Status: types.Present,
		},
		&types.Polygon{
			P:      []types.Vec2{{3.14, -1.678}, {7.1, -5.234}, {23.1, 9.34}},
			Status: types.Present,
		},
		&types.Polygon{Status: types.Null},
	})
}
