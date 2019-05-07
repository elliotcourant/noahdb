package frunk

import (
	"encoding/json"
	"time"
)

// commandType are commands that affect the state of the cluster, and must go through Raft.
type commandType int

const (
	execute        commandType = iota // Commands which modify the database.
	query                             // Commands which query the database.
	metadataSet                       // Commands which sets Store metadata.
	metadataDelete                    // Commands which deletes Store metadata.
	connect                           // Commands which create a database connection.
	disconnect                        // Commands which disconnect from the database.
	kvSet                             // Commands which sets Key-Value data.
	kvDelete                          // Commands which deletes Key-Value data.
)

type command struct {
	Typ commandType     `json:"typ,omitempty"`
	Sub json.RawMessage `json:"sub,omitempty"`
}

func newCommand(t commandType, d interface{}) (*command, error) {
	b, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}
	return &command{
		Typ: t,
		Sub: b,
	}, nil
}

func newMetadataSetCommand(id string, md map[string]string) (*command, error) {
	m := metadataSetSub{
		RaftID: id,
		Data:   md,
	}
	return newCommand(metadataSet, m)
}

func newSetCommand(set map[string]string) (*command, error) {
	k := keyValueSetSub{
		Data:   set,
	}
	return newCommand(kvSet, k)
}

func newDeleteCommand(keys []string) (*command, error) {
	k := keyValueDeleteSub{
		Keys:   keys,
	}
	return newCommand(kvDelete, k)
}

// databaseSub is a command sub which involves interaction with the database.
type databaseSub struct {
	ConnID  uint64   `json:"conn_id,omitempty"`
	Atomic  bool     `json:"atomic,omitempty"`
	Queries []string `json:"queries,omitempty"`
	Timings bool     `json:"timings,omitempty"`
}

type metadataSetSub struct {
	RaftID string            `json:"raft_id,omitempty"`
	Data   map[string]string `json:"data,omitempty"`
}

type connectionSub struct {
	ConnID      uint64        `json:"conn_id,omitempty"`
	IdleTimeout time.Duration `json:"idle_timeout,omitempty"`
	TxTimeout   time.Duration `json:"tx_timeout,omitempty"`
}

type keyValueSetSub struct {
	Data   map[string]string `json:"data,omitempty"`
}

type keyValueDeleteSub struct {
	Keys   []string `json:"keys,omitempty"`
}
