package engine

import (
	"github.com/elliotcourant/noahdb/pkg/transport"
	"github.com/hashicorp/raft"
)

// Config allows custom options to be passed when initializing the engine core.
type Config struct {
	DataDirectory         string
	JoinAddresses         []raft.Server
	Transport             transport.TransportWrapper
	LocalPostgresAddress  string
	LocalPostgresPort     int
	LocalPostgresUsername string
	LocalPostgresPassword string
	StartPool             bool
	AutoJoin              bool
}
