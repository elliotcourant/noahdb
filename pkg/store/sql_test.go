package store_test

import (
	"github.com/elliotcourant/noahdb/pkg/store"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestSqlGeneric(t *testing.T) {
	tmpDir, _ := ioutil.TempDir("", "store_test")
	defer os.RemoveAll(tmpDir)
	store1, err := store.CreateStore(tmpDir, ":", "")
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
	// Simple way to ensure there is a leader.
	time.Sleep(5 * time.Second)

	_, err = store1.Exec("CREATE TABLE foo (id BIGINT PRIMARY KEY);")
	assert.NoError(t, err)

	_, err = store1.Exec("CREATE TABLE foo2 (id BIGINT PRIMARY KEY);")
	assert.NoError(t, err)

	store1.Close()
	store1 = nil
	time.Sleep(5 * time.Second)
	store1, err = store.CreateStore(tmpDir, ":", "")
	defer store1.Close()

}
