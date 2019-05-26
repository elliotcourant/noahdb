package core_test

import (
	"github.com/elliotcourant/noahdb/testutils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDataNodeContext_GetRandomDataNode(t *testing.T) {
	colony, cleanup := testutils.NewTestColony()
	defer cleanup()
	t.Run("get random node", func(t *testing.T) {
		dataNode, err := colony.DataNodes().GetRandomDataNodeShardID()
		assert.NoError(t, err)
		assert.NotEmpty(t, dataNode)
	})
}
