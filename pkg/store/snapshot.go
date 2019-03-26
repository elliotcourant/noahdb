package store

import (
	"github.com/readystock/raft"
)

type snapshot struct {
	store []byte
}

func (s *snapshot) Persist(sink raft.SnapshotSink) (err error) {
	if err = func() error {
		if _, err := sink.Write(s.store); err != nil {
			return err
		}
		return sink.Close()
	}(); err != nil {
		sink.Cancel()
	}
	return err
}

func (s *snapshot) Release() {

}
