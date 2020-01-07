package engine

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetPgDatabaseName(t *testing.T) {
	dataNodeShardID := uint64(1234)
	database := getPgDatabaseName(dataNodeShardID)
	assert.Equal(t, "noahdb_shard_1234", database)
}
