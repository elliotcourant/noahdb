package frunk

import (
	"encoding/json"
)

// Get returns a byte array for a given key from BoltDB.
func (s *Store) Get(key []byte) ([]byte, error) {
	return s.boltStore.Get(key)
}

// Set updates the value of a key-value pair.
func (s *Store) Set(key, value []byte) error {
	if s.IsLeader() {
		cmd, err := newSetCommand(map[string]string{
			string(key): string(value),
		})
		if err != nil {
			return err
		}
		b, err := json.Marshal(cmd)
		if err != nil {
			return err
		}
		return s.raft.Apply(b, applyTimeout).Error()
	}
	return ErrNotLeader
}

// Query allows read-only queries to be issued to the database.
func (s *Store) Query(query string) (*QueryResponse, error) {
	queryRequest := &QueryRequest{
		Queries: []string{
			query,
		},
		Timings: false,
		Atomic:  false,
		Lvl:     None,
	}
	return s.QueryEx(queryRequest)
}

// Exec executes a write-only query against the database.
func (s *Store) Exec(query string) (*ExecuteResponse, error) {
	executeRequest := &ExecuteRequest{
		Queries: []string{
			query,
		},
		Timings: false,
		Atomic:  true,
	}
	return s.ExecuteEx(executeRequest)
}
