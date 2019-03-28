package types_test

import (
	"testing"

	"github.com/elliotcourant/noahdb/pkg/types"
	"github.com/elliotcourant/noahdb/pkg/types/testutil"
)

func TestTIDTranscode(t *testing.T) {
	testutil.TestSuccessfulTranscode(t, "tid", []interface{}{
		&types.TID{BlockNumber: 42, OffsetNumber: 43, Status: types.Present},
		&types.TID{BlockNumber: 4294967295, OffsetNumber: 65535, Status: types.Present},
		&types.TID{Status: types.Null},
	})
}
