package types_test

import (
	"testing"

	"github.com/elliotcourant/noahdb/pkg/types"
	"github.com/elliotcourant/noahdb/pkg/types/testutil"
)

func TestPathTranscode(t *testing.T) {
	testutil.TestSuccessfulTranscode(t, "path", []interface{}{
		&types.Path{
			P:      []types.Vec2{{3.14, 1.678901234}, {7.1, 5.234}},
			Closed: false,
			Status: types.Present,
		},
		&types.Path{
			P:      []types.Vec2{{3.14, 1.678}, {7.1, 5.234}, {23.1, 9.34}},
			Closed: true,
			Status: types.Present,
		},
		&types.Path{
			P:      []types.Vec2{{7.1, 1.678}, {-13.14, -5.234}},
			Closed: true,
			Status: types.Present,
		},
		&types.Path{Status: types.Null},
	})
}
