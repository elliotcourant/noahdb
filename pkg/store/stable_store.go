package store

import (
	"github.com/dgraph-io/badger"
	"github.com/kataras/golog"
)

type stableStore Store

func (stable *stableStore) Set(key, val []byte) error {
	return stable.badger.Update(func(txn *badger.Txn) error {
		golog.Debugf("Setting Key: %s To Value: %s", string(key), string(val))
		return txn.Set(key, val)
	})
}

func (stable *stableStore) Get(key []byte) (val []byte, err error) {
	err = stable.badger.View(func(txn *badger.Txn) (err error) {
		if item, err := txn.Get(key); err != nil {
			return err
		} else {
			val := make([]byte, 0)
			val, err = item.ValueCopy(val)
			return err
		}
	})
	if err != nil && err.Error() == "Key not found" {
		return make([]byte, 0), nil
	}
	return val, err
}

func (stable *stableStore) SetUint64(key []byte, val uint64) error {
	return stable.Set(key, Uint64ToBytes(val))
}

func (stable *stableStore) GetUint64(key []byte) (uint64, error) {
	if val, err := stable.Get(key); err != nil {
		return 0, err
	} else {
		return BytesToUint64(val), nil
	}
}
