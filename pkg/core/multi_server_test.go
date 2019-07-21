package core_test

import (
	"github.com/elliotcourant/noahdb/pkg/core"
	"github.com/elliotcourant/noahdb/testutils"
	"github.com/elliotcourant/timber"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestMultiServer_LeaderFailure(t *testing.T) {
	initialLeader, initialCleanup := testutils.NewTestColony(t)

	numberOfFollowers := 7

	followers := make([]core.Colony, numberOfFollowers)
	cleanups := make([]func(), numberOfFollowers)
	for i := 0; i < numberOfFollowers; i++ {
		followerColony, followerCleanup := testutils.NewTestColony(t, initialLeader.Addr().String())
		cleanups[i] = followerCleanup
		followers[i] = followerColony
	}

	defer func(cleanups []func()) {
		for _, cleanup := range cleanups {
			cleanup()
		}
	}(cleanups)
	time.Sleep(10 * time.Second)

	initialCleanup()

	time.Sleep(10 * time.Second)

	foundNewLeader := false
	for _, follower := range followers {
		if follower.IsLeader() {
			timber.Infof("Found new leader!")
			foundNewLeader = true
			break
		}
	}
	assert.True(t, foundNewLeader)
}

func TestMultiServerWrite(t *testing.T) {
	t.Skip()
	t.Run("create a new schema", func(t *testing.T) {
		colony1, cleanup1 := testutils.NewTestColony(t)
		defer cleanup1()

		colony2, cleanup2 := testutils.NewTestColony(t, colony1.Addr().String())
		defer cleanup2()

		// Make sure the first colony is the leader.
		assert.True(t, colony1.IsLeader())
		assert.False(t, colony2.IsLeader())

		// Try to create a schema from a follower node.
		followerSchemaName := "test"
		followerSchema, err := colony2.Schema().NewSchema(followerSchemaName)
		assert.NoError(t, err)
		assert.NotEmpty(t, followerSchema)
		assert.True(t, followerSchema.SchemaID > 0)
		assert.Equal(t, followerSchemaName, followerSchema.SchemaName)
	})
}
