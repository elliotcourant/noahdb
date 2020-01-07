package engine

import (
	"fmt"
)

// getPgDatabaseName will return the database name string that will be used to connect to the actual
// shard in PostgreSQL. This is simply a string with the dataNodeShardId suffix.
func getPgDatabaseName(dataNodeShardId uint64) string {
	return fmt.Sprintf("noahdb_shard_%d", dataNodeShardId)
}
