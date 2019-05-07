package frunk

import (
	"encoding/json"
)

func (s *Store) Get(key []byte) ([]byte, error) {
	return s.boltStore.Get(key)
}

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
