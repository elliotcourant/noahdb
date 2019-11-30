package engine_test

import (
	"github.com/elliotcourant/meles"
	"github.com/elliotcourant/mellivora"
	"github.com/elliotcourant/noahdb/pkg/engine"
	"github.com/elliotcourant/timber"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"math/rand"
	"net"
	"os"
	"sync"
	"testing"
)

func TestMain(m *testing.M) {
	res := m.Run()
	os.Exit(res)
}

type TestCluster []engine.Core

func (tc TestCluster) Begin(t *testing.T) engine.Transaction {
	return tc.BeginOn(t, rand.Int())
}

func (tc TestCluster) BeginOn(t *testing.T, node int) engine.Transaction {
	txn, err := tc[node%len(tc)].Begin()
	if !assert.NoError(t, err) {
		panic(err)
	}

	return txn
}

func NewTestCoreCluster(t *testing.T, numberOfPeers int) (TestCluster, func()) {
	peers := make([]string, numberOfPeers)
	listeners := make([]net.Listener, numberOfPeers)
	dirs := make([]string, numberOfPeers)
	for i := 0; i < numberOfPeers; i++ {
		listener, err := net.Listen("tcp", ":")
		if !assert.NoError(t, err) {
			panic(err)
		}

		tmpDir, err := ioutil.TempDir("", "noahdb-core")
		if !assert.NoError(t, err) {
			panic(err)
		}

		listeners[i] = listener
		dirs[i] = tmpDir
		peers[i] = listener.Addr().String()
	}

	baseLogger := timber.With(timber.Keys{
		"test": t.Name(),
	})

	wg := sync.WaitGroup{}
	wg.Add(numberOfPeers)

	cluster := make(TestCluster, numberOfPeers)
	for i := 0; i < numberOfPeers; i++ {
		logger := baseLogger.With(timber.Keys{
			"node": i,
		})

		store, err := meles.NewStore(listeners[i], logger, meles.Options{
			Directory: dirs[i],
			Peers:     peers,
		})
		if !assert.NoError(t, err) {
			panic(err)
		}

		go func() {
			defer wg.Done()
			if err := store.Start(); !assert.NoError(t, err) {
				panic(err)
			}
		}()

		db := mellivora.NewDatabase(store, logger)

		cluster[i] = engine.NewCore(store, db)
	}
	wg.Wait()

	return cluster, func() {
		for _, node := range cluster {
			node.Close()
		}
		for _, listener := range listeners {
			listener.Close()
		}
		for _, dir := range dirs {
			os.RemoveAll(dir)
		}
	}
}
