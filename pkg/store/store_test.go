package store_test

import (
	"fmt"
	"github.com/ahmetb/go-linq"
	"github.com/elliotcourant/noahdb/pkg/store"
	"github.com/kataras/golog"
	"io/ioutil"
	"os"
	"reflect"
	"sort"
	"sync"
	"testing"
	"time"
)

func TestCreateStore(t *testing.T) {
	tmpDir, _ := ioutil.TempDir("", "store_test")
	defer os.RemoveAll(tmpDir)
	store1, err := store.CreateStore(tmpDir, "127.0.0.1:0", "")
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
	defer store1.Close()
	// Simple way to ensure there is a leader.
	time.Sleep(5 * time.Second)
	err = store1.Set([]byte("test"), []byte("value"))
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}

	val, err := store1.Get([]byte("test"))
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}

	if string(val) != "value" {
		t.Error("value did not match")
		t.Fail()
		return
	}
}

func TestCreateStoreSeveralServers(t *testing.T) {
	serverCount := 6
	startingPort := 7543
	golog.Infof("starting %d server(s) for testing", serverCount)
	tmpDirs := make([]string, serverCount)
	for i := 0; i < serverCount; i++ {
		if tmpDir, err := ioutil.TempDir("", fmt.Sprintf("store_test_%d", i)); err != nil {
			panic(err)
		} else {
			tmpDirs[i] = tmpDir
		}
	}
	defer func() {
		for _, tmpDir := range tmpDirs {
			if err := os.RemoveAll(tmpDir); err != nil {
				panic(err)
			}
		}
	}()

	stores := make(chan *store.Store, serverCount)
	var wg sync.WaitGroup
	wg.Add(serverCount)
	for i := 0; i < serverCount; i++ {
		listenPort := fmt.Sprintf(":%d", startingPort)
		joinPort := ""
		if i == 0 {
			golog.Infof("starting node [%d] raft port [%s]", i, listenPort)
		} else {
			joinPort = fmt.Sprintf(":%d", startingPort-(2*i))
			golog.Infof("starting node [%d] raft port [%s] joining [%s]", i, listenPort, joinPort)
		}

		startingPort += 2 // Increment the starting port for the next iteration

		tmpDir := tmpDirs[i]
		go func(index int, tmpDir, listenPort, joinPort string) {
			defer wg.Done()
			store, err := store.CreateStore(tmpDir, listenPort, joinPort)
			if err != nil {
				panic(err)
			}
			time.Sleep(5 * time.Second)
			stores <- store
			golog.Infof("finished starting node [%d]", index)
		}(i, tmpDir, listenPort, joinPort)
		time.Sleep(10 * time.Second)
	}
	wg.Wait()
	close(stores)

	tupleSync := sync.Mutex{}
	tupleIndex := 0
	tuples := make([]store.Tuple, serverCount)
	addTuple := func(tuple store.Tuple) {
		tupleSync.Lock()
		defer tupleSync.Unlock()
		tuples[tupleIndex] = tuple
		tupleIndex++
	}

	servers := make(chan *store.Store, serverCount)
	wg = sync.WaitGroup{}
	wg.Add(serverCount)
	for server := range stores {
		go func(str *store.Store) {
			defer wg.Done()
			id := str.NodeID()
			if str.IsLeader() {
				golog.Infof("node [%d] is leader", id)
			} else {
				golog.Infof("node [%d] is follower", id)
			}

			tuple := store.Tuple{
				Key:   []byte(fmt.Sprintf("key_%d", id)),
				Value: []byte(fmt.Sprintf("value_%d", id)),
			}
			golog.Infof("[%d] starting to set [%s]", str.NodeID(), string(tuple.Key))

			if err := str.Set(tuple.Key, tuple.Value); err != nil {
				panic(err)
			}

			golog.Infof("[%d] finished setting [%s]", str.NodeID(), string(tuple.Key))
			addTuple(tuple)
			servers <- str
			time.Sleep(5 * time.Second)
		}(server)
	}

	fmt.Printf("WAITING FOR ALL SERVERS TO FINISH.\n")
	wg.Wait()
	fmt.Printf("ALL %d SERVERS HAVE FINISHED, NOW READING.\n", len(servers))
	close(servers)
	wg = sync.WaitGroup{}
	wg.Add(serverCount)
	for server := range servers {
		go func(store *store.Store) {
			defer wg.Done()
			id := store.NodeID()
			for _, tuple := range tuples {
				if value, err := store.Get(tuple.Key); err != nil {
					panic(err)
				} else if !reflect.DeepEqual(value, tuple.Value) {
					panic(fmt.Sprintf("Tuple for key [%s] does not match on node [%d] expected [%s] found [%s]", string(tuple.Key), id, string(value), string(tuple.Value)))
				}
			}
		}(server)
	}
	fmt.Printf("WAITING FOR ALL SERVERS TO FINISH READS.\n")
	wg.Wait()
	fmt.Printf("ALL SERVERS HAVE FINISHED\n")
}

func TestGetPrefix(t *testing.T) {
	tmpDir, _ := ioutil.TempDir("", "store_test")
	defer os.RemoveAll(tmpDir)
	store1, err := store.CreateStore(tmpDir, ":6802", "")
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
	defer store1.Close()
	// Simple way to ensure there is a leader.
	time.Sleep(5 * time.Second)
	err = store1.Set([]byte("/test"), []byte("value"))
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}

	val, err := store1.GetPrefix([]byte("/"))
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
	if len(val) == 0 {
		t.Error("no values found")
		t.Fail()
		return
	}
	for _, kv := range val {
		golog.Debugf("Key: %s Value: %s", string(kv.Key), string(kv.Value))
	}
}

func TestSequence(t *testing.T) {
	tmpDir, _ := ioutil.TempDir("", "store_test")
	defer os.RemoveAll(tmpDir)
	store1, err := store.CreateStore(tmpDir, ":6702", "")
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
	defer store1.Close()
	// Simple way to ensure there is a leader.
	time.Sleep(5 * time.Second)
	numberOfIds := 10000
	Ids := make([]int, 0)
	for i := 0; i < numberOfIds; i++ {
		id, err := store1.NextSequenceValueById("public.users.user_id")
		if err != nil {
			t.Error(err)
			t.Fail()
			return
		}
		Ids = append(Ids, int(*id))
		// golog.Infof("New user_id: %d", *id)
	}
	sort.Ints(Ids)
	if len(Ids) != numberOfIds {
		t.Error("number of ids do not match")
		t.Fail()
		return
	}

}

func TestSequenceMulti(t *testing.T) {
	tmpDir, _ := ioutil.TempDir("", "store_test")
	defer os.RemoveAll(tmpDir)
	tmpDir2, _ := ioutil.TempDir("", "store_test2")
	defer os.RemoveAll(tmpDir2)
	store1, err := store.CreateStore(tmpDir, ":6502", "")
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
	defer store1.Close()
	// Simple way to ensure there is a leader.
	time.Sleep(5 * time.Second)

	store2, err := store.CreateStore(tmpDir2, ":6503", ":6502")
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
	defer store2.Close()
	// store1.Join(store2.NodeID(), ":6546", ":6501")
	time.Sleep(5 * time.Second)

	numberOfIds := 1000000
	Ids := make([]int, 0)
	for i := 0; i < numberOfIds; i++ {
		switch i % 2 {
		case 0:
			id, err := store2.NextSequenceValueById("public.users.user_id")
			if err != nil {
				panic(err)
				t.Fail()
				return
			}
			Ids = append(Ids, int(*id))
			// golog.Infof("New user_id on node 2: %d", *id)
		default:
			id, err := store1.NextSequenceValueById("public.users.user_id")
			if err != nil {
				panic(err)
				t.Fail()
				return
			}
			Ids = append(Ids, int(*id))
			// golog.Infof("New user_id on node 1: %d", *id)
		}

	}
	sort.Ints(Ids)
	if len(Ids) != numberOfIds {
		t.Error("number of ids do not match")
		t.Fail()
		return
	}
	linq.From(Ids).Distinct().ToSlice(&Ids)
	if len(Ids) != numberOfIds {
		t.Error("distinct number of ids do not match")
		t.Fail()
		return
	}
}

func TestCreateStoreWithClose(t *testing.T) {
	tmpDir, _ := ioutil.TempDir("", "store_test")
	defer os.RemoveAll(tmpDir)
	store1, err := store.CreateStore(tmpDir, ":6602", "")
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
	// Simple way to ensure there is a leader.
	time.Sleep(5 * time.Second)
	err = store1.Set([]byte("test"), []byte("value"))
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}

	val, err := store1.Get([]byte("test"))
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}

	if string(val) != "value" {
		t.Error("value did not match")
		t.Fail()
		return
	}

	golog.Warnf("shutting down lone node and restarting it")
	store1.Close()
	store1 = nil
	time.Sleep(5 * time.Second)
	store1, err = store.CreateStore(tmpDir, ":6602", "")
	defer store1.Close()
	// Simple way to ensure there is a leader.
	time.Sleep(5 * time.Second)
	err = store1.Set([]byte("test1"), []byte("value"))
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}

	val, err = store1.Get([]byte("test1"))
	if err != nil {
		t.Error(err)
		t.Fail()
		return
	}
}
