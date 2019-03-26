package store

import (
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"testing"
	"time"
)

func TestSequenceServer_GetSequenceChunk(t *testing.T) {
	tmpDir, _ := ioutil.TempDir("", "store_test")
	defer os.RemoveAll(tmpDir)
	store1, err := CreateStore(tmpDir, "127.0.0.1:0", "")
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
	// Simple way to ensure there is a leader.
	time.Sleep(5 * time.Second)
	sequence := clusterServer{*store1}
	response1, err := sequence.getSequenceChunk("test")
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
	response2, err := sequence.getSequenceChunk("test")
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
	chunk := &SequenceChunk{
		current: response1,
		next:    response2,
		index:   1,
		sync:    new(sync.Mutex),
	}

	for i := 0; i < 10; i++ {
		val, err := chunk.Next()
		if err != nil {
			t.Error(err)
			t.Fail()
			return
		}
		fmt.Printf("Sequence Value: %d\n", *val)
	}
	response, err := sequence.getSequenceChunk("test")
	chunk.current = response
	chunk.index = 1
	val, err := chunk.Next()
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
	fmt.Printf("Sequence Value: %d\n", *val)
}
