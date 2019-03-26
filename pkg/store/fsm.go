package store

import (
	"bytes"
	"database/sql"
	"encoding/hex"
	"fmt"
	"github.com/dgraph-io/badger"
	_ "github.com/elliotcourant/noahdb/pkg/drivers/sqlite"
	"github.com/golang/protobuf/proto"
	"github.com/readystock/golog"
	"github.com/readystock/raft"
	"io"
	"time"
)

type fsm Store

func (f *fsm) Apply(l *raft.Log) interface{} {
	var c Command
	if err := proto.Unmarshal(l.Data, &c); err != nil {
		golog.Fatalf("failed to unmarshal command: %s. %s", err.Error(), hex.Dump(l.Data))
		return err
	}
	r := CommandResponse{
		Timestamp: c.Timestamp,
		Operation: c.Operation,
	}
	golog.Verbosef("[%d] FSM Receive Delay [%s]", f.nodeId, time.Since(time.Unix(0, int64(c.Timestamp))))
	if err := func() error {
		switch c.OperationType {
		case OperationType_SQL:
			golog.Verbosef("[%d] FSM Received Query: [%s]", f.nodeId, string(c.Query))
			_, err := f.sqlstore.Exec(string(c.Query))
			return err
		case OperationType_KV:
			switch c.Operation {
			case Operation_SET:
				return f.applySet(c.Key, c.Value)
			case Operation_DELETE:
				return f.applyDelete(c.Key)
			default:
				return fmt.Errorf("unsupported command operation: %s", c.Operation)
			}
		default:
			return fmt.Errorf("unsupported command operation type: %s", c.OperationType)
		}
	}(); err != nil {
		r.ErrorMessage = err.Error()
		r.IsSuccess = false
	} else {
		r.IsSuccess = true
		r.AppliedTimestamp = uint64(time.Now().UnixNano())
	}
	return r
}

// Restore stores the key-value store to a previous state.
func (f *fsm) Restore(rc io.ReadCloser) error {
	if err := f.badger.Load(rc); err != nil {
		return err
	}
	logEntries, err := f.GetPrefix(logsPrefix)
	if err != nil {
		return err
	}
	f.sqlstore, err = sql.Open("sqlite3", ":memory:")
	if err != nil {
		return err
	}
	for _, entry := range logEntries {
		logEntry := raft.Log{}
		if err := proto.Unmarshal(entry.Value, &logEntry); err != nil {
			return err
		}
		command := Command{}
		if err := proto.Unmarshal(logEntry.Data, &command); err != nil {
			continue
		}
		if command.OperationType == OperationType_SQL && len(command.Query) > 0 {
			golog.Verbosef("Replaying SQL: [%s]", string(command.Query))
			if _, err := f.sqlstore.Exec(string(command.Query)); err != nil {
				return err
			}
		}
	}
	return nil
}

func (f *fsm) GetPrefix(prefix []byte) (values []KeyValue, err error) {
	values = make([]KeyValue, 0)
	err = f.badger.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			keyBytes := make([]byte, 0)
			valueBytes := make([]byte, 0)
			keyBytes = item.KeyCopy(keyBytes)
			if valueBytes, err = item.ValueCopy(valueBytes); err != nil {
				return err
			} else {
				values = append(values, KeyValue{Key: keyBytes, Value: valueBytes})
			}
		}
		return nil
	})
	return values, err
}

// Snapshot returns a snapshot of the key-value store.
func (f *fsm) Snapshot() (raft.FSMSnapshot, error) {
	w := &bytes.Buffer{}
	f.badger.Backup(w, 0)
	return &snapshot{
		store: w.Bytes(),
	}, nil
}

func (f *fsm) applySet(key, value []byte) error {
	return f.badger.Update(func(txn *badger.Txn) error {
		golog.Verbosef("[%d] FSM Setting Key: %s", f.nodeId, string(key))
		return txn.Set(key, value)
	})
}

func (f *fsm) applyDelete(key []byte) error {
	return f.badger.Update(func(txn *badger.Txn) error {
		golog.Verbosef("[%d] Deleting Key: %s", f.nodeId, string(key))
		return txn.Delete(key)
	})
}
