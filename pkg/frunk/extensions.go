package frunk

import (
	"encoding/json"
	"fmt"
	"github.com/elliotcourant/timber"
	"time"
)

// Get returns a byte array for a given key from BoltDB.
func (s *Store) Get(key []byte) ([]byte, error) {
	val, err := s.boltStore.Get(key)
	if err != nil && err.Error() != "not found" {
		return nil, err
	}
	return val, nil
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
	startTimestamp := time.Now()
	defer func() {
		timber.New().SetDepth(2).Tracef("[%s] %s", time.Since(startTimestamp), query)
	}()
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
	return s.ExecEx(query, true)
}

// ExecEx executes a write-only query against the database.
// With the option of being atomic or not.
func (s *Store) ExecEx(query string, atomic bool) (*ExecuteResponse, error) {
	startTimestamp := time.Now()
	defer func() {
		timber.New().SetDepth(2).Tracef("[%s] %s", time.Since(startTimestamp), query)
	}()
	executeRequest := &ExecuteRequest{
		Queries: []string{
			query,
		},
		Timings: false,
		Atomic:  atomic,
	}
	result, err := s.ExecuteEx(executeRequest)
	if err != nil {
		return nil, err
	}
	if len(result.Results) > 0 {
		for _, resultSet := range result.Results {
			if resultSet.Error != "" {
				return nil, fmt.Errorf(resultSet.Error)
			}
		}
	}
	return result, err
}
