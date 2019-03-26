package store

import (
	"encoding/hex"
	"fmt"
	"github.com/dgraph-io/badger"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"os"
	"testing"
)

func newLogStore() (*logStore, error) {
	store := Store{}
	randomName, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	tmpDir, _ := ioutil.TempDir("", randomName.String())
	defer os.RemoveAll(tmpDir)
	opts := badger.DefaultOptions
	opts.Dir = tmpDir
	opts.ValueDir = tmpDir
	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}
	store.badger = db
	ls := logStore(store)
	return &ls, nil
}

func TestLogStore_FirstIndex(t *testing.T) {
	ls, err := newLogStore()
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}

	if index, err := ls.FirstIndex(); err != nil {
		t.Error(err)
		t.Fail()
		return
	} else {
		if index != 0 {
			t.Error("first index should be 0")
			t.Fail()
			return
		}
	}
}

func TestLogStore_LastIndex(t *testing.T) {
	ls, err := newLogStore()
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}

	if index, err := ls.LastIndex(); err != nil {
		t.Error(err)
		t.Fail()
		return
	} else {
		if index != 0 {
			t.Error("last index should be 0")
			t.Fail()
			return
		}
	}
}

func TestLogStore_GetKeyForIndex(t *testing.T) {
	uints := []uint64{1, 2, 3, 15831904231, 35183541, 489156156156}
	for _, u := range uints {
		key := getKeyForIndex(u)
		index := getIndexForKey(key)
		fmt.Print(hex.Dump(key))
		if index != u {
			t.Error("decoded index does not match input index")
			t.Fail()
			return
		}
	}
}
