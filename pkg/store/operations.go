package store

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/dgraph-io/badger"
	"github.com/golang/protobuf/proto"
	"github.com/readystock/golog"
	"github.com/readystock/raft"
	"time"
)

func (store *Store) GetPrefixWithPredicate(prefix []byte, predicate func(kv KeyValue) (bool, error)) (kv KeyValue, found bool, err error) {
	return kv, found, store.badger.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			keyBytes := make([]byte, 0)
			keyBytes = item.KeyCopy(keyBytes)
			valueBytes := make([]byte, 0)
			valueBytes, err = item.ValueCopy(valueBytes)
			if err != nil {
				return err
			}
			val := KeyValue{Key: keyBytes, Value: valueBytes}
			if f, err := predicate(val); err != nil {
				return err
			} else if f {
				kv = val
				found = true
				return nil
			}
		}
		return nil
	})
}

func (store *Store) GetPrefix(prefix []byte) (values []KeyValue, err error) {
	values = make([]KeyValue, 0)
	err = store.badger.View(func(txn *badger.Txn) error {
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

func (store *Store) GetKeyOnlyPrefix(prefix []byte) (keys [][]byte, err error) {
	keys = make([][]byte, 0)
	err = store.badger.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.IteratorOptions{
			PrefetchValues: false,
			Reverse:        false,
			AllVersions:    false,
			Prefix:         prefix,
		})
		defer it.Close()
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			key := make([]byte, 0)
			key = item.KeyCopy(key)
			keys = append(keys, key)
		}
		return nil
	})
	return keys, err
}

func (store *Store) Get(key []byte) (value []byte, err error) {
	// isSet := true
	err = store.badger.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			if err.Error() != "Key not found" {
				return err
			} else {
				// If the key is not found normally we want to just return an empty byte array.
				// but only if we are the leader, if we are not the leader then we want to send
				// a request to the leader.
				// isSet = false
				value = make([]byte, 0)
				return nil
			}
		}
		value, err = item.ValueCopy(value)
		return err
	})
	// if !isSet {
	//     value, err = store.clusterClient.Get(key)
	// }
	return value, err
}

func (store *Store) Set(key, value []byte) (err error) {
	c := &Command{OperationType: OperationType_KV, Operation: Operation_SET, Key: key, Value: value, Timestamp: uint64(time.Now().UnixNano())}
	if store.raft.State() != raft.Leader {
		if store.raft.Leader() == "" {
			return errors.New("no leader in cluster")
		}
		if _, err := store.clusterClient.sendCommand(c); err != nil {
			return err
		}
		return nil
	}
	b, err := proto.Marshal(c)
	if err != nil {
		return err
	}
	r := store.raft.Apply(b, raftTimeout)
	if err := r.Error(); err != nil {
		return err
	} else if resp, ok := r.Response().(CommandResponse); ok {
		golog.Debugf("[%d] Delay Total [%s] Response [%s]", store.nodeId, time.Since(time.Unix(0, int64(resp.Timestamp))), time.Since(time.Unix(0, int64(resp.AppliedTimestamp))))
	}
	return nil
}

func (store *Store) ChangeKey(currentKey []byte, newKey []byte) error {
	val, err := store.Get(currentKey)
	if err != nil {
		return err
	}

	if err := store.Set(newKey, val); err != nil {
		return err
	}

	return store.Delete(currentKey)
}

func (store *Store) Delete(key []byte) (err error) {
	c := &Command{OperationType: OperationType_KV, Operation: Operation_DELETE, Key: key, Value: nil, Timestamp: uint64(time.Now().UnixNano())}
	if store.raft.State() != raft.Leader {
		if _, err := store.clusterClient.sendCommand(c); err != nil {
			return err
		}
		return nil
	}
	b, err := proto.Marshal(c)
	if err != nil {
		return err
	}
	r := store.raft.Apply(b, raftTimeout)
	return r.Error()
}

func (store *Store) Exec(query string) (sql.Result, error) {
	c := &Command{OperationType: OperationType_SQL, Timestamp: uint64(time.Now().UnixNano()), Query: []byte(query)}
	if store.raft.State() != raft.Leader {
		if _, err := store.clusterClient.sendCommand(c); err != nil {
			return nil, err
		}
		return nil, nil
	}
	b, err := proto.Marshal(c)
	if err != nil {
		return nil, err
	}
	r := store.raft.Apply(b, raftTimeout)
	if err := r.Error(); err != nil {
		return nil, err
	} else if resp, ok := r.Response().(CommandResponse); ok {
		golog.Debugf("[%d] Delay Total [%s] Response [%s]", store.nodeId, time.Since(time.Unix(0, int64(resp.Timestamp))), time.Since(time.Unix(0, int64(resp.AppliedTimestamp))))
		if !resp.IsSuccess {
			return nil, fmt.Errorf(resp.ErrorMessage)
		}
	}
	return nil, r.Error()
}

func (store *Store) Query(query string) (*sql.Rows, error) {
	return store.sqlstore.Query(query)
}
