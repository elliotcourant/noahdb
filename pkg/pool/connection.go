package pool

import (
	"github.com/elliotcourant/noahdb/pkg/pgproto"
)

type (
	Connection interface {
		pgproto.Frontend
		Close() error
		Release()
		DataNodeID() uint64
		ShardID() uint64
		DataNodeShardID() uint64
		IsRoot() bool
	}
)
