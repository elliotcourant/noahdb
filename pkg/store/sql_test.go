package store_test

import (
	"fmt"
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

func TestSqlInserts(t *testing.T) {
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

	store1.Close()
	store1 = nil
	time.Sleep(5 * time.Second)
	store1, err = store.CreateStore(tmpDir, ":9001", "")
	assert.NoError(t, err)
	time.Sleep(5 * time.Second)
	for i := 0; i < 100; i++ {
		_, err := store1.Exec(fmt.Sprintf("INSERT INTO foo (id) VALUES(%d);", i))
		assert.NoError(t, err)
	}

	tmpDir2, _ := ioutil.TempDir("", "store_test2")
	defer os.RemoveAll(tmpDir2)
	store2, err := store.CreateStore(tmpDir2, ":9002", ":9001")
	time.Sleep(5 * time.Second)
	rows, err := store2.Query("SELECT COUNT(id) FROM foo;")
	assert.NoError(t, err)
	rows.Next()
	countVal := 0
	rows.Scan(&countVal)
	assert.Equal(t, 100, countVal)
	time.Sleep(10 * time.Second)
	store1.Close()
	store2.Close()
}
