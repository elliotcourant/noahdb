package core

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDataNodeContext_GetRandomDataNode(t *testing.T) {
	colony, cleanup := newTestColony()
	defer cleanup()
	t.Run("get random node", func(t *testing.T) {
		dataNode, err := colony.DataNodes().(*dataNodeContext).GetRandomDataNodeShardID()
		assert.NoError(t, err)
		assert.Empty(t, dataNode)
	})
}
