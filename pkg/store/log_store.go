package store

import (
	"github.com/dgraph-io/badger"
	"github.com/gogo/protobuf/proto"
	"github.com/readystock/raft"
)

type logStore Store

var (
	logsPrefix = []byte("/_logs_/")
)

func (log *logStore) FirstIndex() (val uint64, err error) {
	return log.index(false)
}

func (log *logStore) LastIndex() (uint64, error) {
	return log.index(true)
}

func (log *logStore) GetLog(index uint64, raftLog *raft.Log) error {
	return log.badger.View(func(txn *badger.Txn) (err error) {
		item, err := txn.Get(getKeyForIndex(index))
		if err != nil {
			return raft.ErrLogNotFound
		}
		value := make([]byte, 0)
		value, err = item.ValueCopy(value)
		if err != nil {
			return err
		}
		if err = proto.Unmarshal(value, raftLog); err != nil {
			return err
		}
		return nil
	})
}

func (log *logStore) StoreLog(raftLog *raft.Log) error {
	return log.StoreLogs([]*raft.Log{raftLog})
}

func (log *logStore) StoreLogs(raftLogs []*raft.Log) error {
	return log.badger.Update(func(txn *badger.Txn) error {
		for _, log := range raftLogs {
			key := getKeyForIndex(log.Index)
			val, err := proto.Marshal(log)
			if err != nil {
				return err
			}
			err = txn.Set(key, val)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (log *logStore) DeleteRange(min, max uint64) error {
	return log.badger.Update(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.IteratorOptions{
			PrefetchSize: 100,
		})
		defer it.Close()
		minKey := Uint64ToBytes(min)
		keys := make([][]byte, 0)
		// Get all the keys in the range
		for it.Seek(minKey); it.Valid(); it.Next() {
			index := getIndexForKey(it.Item().Key())
			if index > max {
				break
			}
			keys = append(keys, it.Item().Key())
		}
		// Delete all of the keys found
		for _, key := range keys {
			if err := txn.Delete(key); err != nil {
				return err
			}
		}
		return nil
	})
	return nil
}

func getKeyForIndex(index uint64) []byte {
	key := make([]byte, 0)
	key = append(key, logsPrefix...)
	key = append(key, Uint64ToBytes(index)...)
	return key
}

func getIndexForKey(key []byte) uint64 {
	return BytesToUint64(key[len(logsPrefix):])
}

func (log *logStore) index(reverse bool) (val uint64, err error) {
	val = 0
	err = log.badger.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.IteratorOptions{
			PrefetchSize:   1,
			PrefetchValues: false,
			Reverse:        reverse,
		})
		defer it.Close()
		for it.Seek(logsPrefix); it.ValidForPrefix(logsPrefix); it.Next() {
			val = getIndexForKey(it.Item().Key())
			return nil
		}
		return nil
	})
	return val, err
}
