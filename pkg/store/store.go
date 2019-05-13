package store

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/dgraph-io/badger"
	"github.com/kataras/golog"
	"github.com/readystock/raft"
	"google.golang.org/grpc"
	"net"
	"os"
	"sync"
	"time"
)

const (
	retainSnapshotCount = 2
	raftTimeout         = 10 * time.Second
)

type KeyValue struct {
	Key   []byte
	Value []byte
}

type Store struct {
	raft              *raft.Raft
	badger            *badger.DB
	sequenceIds       *badger.Sequence
	chunkMapMutex     *sync.Mutex
	sequenceCacheSync *sync.Mutex
	sequenceChunks    map[string]*SequenceChunk
	sequenceCache     map[string]*Sequence
	clusterClient     *clusterClient
	server            *grpc.Server
	nodeId            uint64
	listen            string
	sqlstore          *sql.DB
}

// Creates and possibly joins a cluster.
func CreateStore(directory string, listen string, joinAddr string) (*Store, error) {
	// Setup Raft configuration.
	config := raft.DefaultConfig()
	store := Store{
		chunkMapMutex:     new(sync.Mutex),
		sequenceCacheSync: new(sync.Mutex),
		sequenceCache:     map[string]*Sequence{},
		sequenceChunks:    map[string]*SequenceChunk{},
		listen:            listen,
	}

	sqlstore, _ := sql.Open("sqlite3", ":memory:")
	store.sqlstore = sqlstore
	if listen == "" {
		listen = ":6543"
	}

	lis, err := net.Listen("tcp", listen)
	if err != nil {
		return nil, err
	}
	store.listen = lis.Addr().String()
	grpcServer := grpc.NewServer()
	transport, err := raft.NewGrpcTransport(grpcServer, lis.Addr().String())
	if err != nil {
		return nil, err
	}

	opts := badger.DefaultOptions
	opts.Dir = directory
	opts.ValueDir = directory
	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}
	store.badger = db
	stable := stableStore(store)
	log := logStore(store)

	nodeId := uint64(0)

	clusterExists := false

	if nodeIdBytes, _ := store.Get(serverIdPath); len(nodeIdBytes) > 0 {
		clusterExists = true
		nodeId = BytesToUint64(nodeIdBytes)
	}

	if joinAddr != "" {
		conn, err := grpc.Dial(joinAddr, grpc.WithInsecure())
		if err != nil {
			return nil, err
		}
		tempClient := &clusterServiceClient{cc: conn}

		if !clusterExists {
			response, err := tempClient.GetNodeID(context.Background(), &GetNodeIdRequest{})
			if err != nil {
				return nil, err
			}
			nodeId = response.NodeID
			if err := stable.Set(serverIdPath, Uint64ToBytes(nodeId)); err != nil {
				return nil, err
			}
		}

		defer func() {
			golog.Debugf("node %d joining cluster at addr %s!", nodeId, joinAddr)
			if _, err := tempClient.Join(context.Background(), &JoinRequest{RaftAddress: listen, ID: nodeId}); err != nil {
				golog.Errorf("could not join `%s` error: %s", listen, err)
			}
		}()
	} else {
		if !clusterExists {
			if err := stable.Set(serverIdPath, Uint64ToBytes(nodeId)); err != nil {
				return nil, err
			}
		}
	}

	config.LocalID = raft.ServerID(nodeId)
	store.nodeId = nodeId
	snapshots, err := raft.NewFileSnapshotStore(directory, retainSnapshotCount, os.Stderr)
	if err != nil {
		return nil, fmt.Errorf("file snapshot store: %s", err)
	}
	// if clusterExists {
	// 	configuration := raft.Configuration{
	// 		Servers: []raft.Server{
	// 			{
	// 				ID:      config.LocalID,
	// 				Address: transport.LocalAddr(),
	// 			},
	// 		},
	// 	}
	// 	err := raft.RecoverCluster(config, (*fsm)(&store), &log, &stable, snapshots, transport, configuration)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("recover raft: %s", err)
	// 	}
	// }
	ra, err := raft.NewRaft(config, (*fsm)(&store), &log, &stable, snapshots, transport)
	if err != nil {
		return nil, fmt.Errorf("new raft: %s", err)
	}
	store.raft = ra
	RegisterClusterServiceServer(grpcServer, &clusterServer{store})
	go grpcServer.Serve(lis)
	if joinAddr == "" && nodeId == 0 && !clusterExists {
		configuration := raft.Configuration{
			Servers: []raft.Server{
				{
					ID:      config.LocalID,
					Address: transport.LocalAddr(),
				},
			},
		}
		f := store.raft.BootstrapCluster(configuration)
		if f.Error() != nil {
			return nil, f.Error()
		}
		time.Sleep(5 * time.Second)
	}

	store.server = grpcServer
	store.clusterClient = &clusterClient{Store: store, sync: new(sync.Mutex)}
	return &store, nil
}

func (store *Store) join(nodeId uint64, addr string) error {
	golog.Debugf("received join request from remote node [%d] at [%s]", nodeId, addr)

	configFuture := store.raft.GetConfiguration()
	if err := configFuture.Error(); err != nil {
		golog.Errorf("failed to get raft configuration: %s", err.Error())
	}

	for _, srv := range configFuture.Configuration().Servers {
		// If a node already exists with either the joining node's ID or address,
		// that node may need to be removed from the config first.
		if srv.ID == raft.ServerID(nodeId) || srv.Address == raft.ServerAddress(addr) {
			// However if *both* the ID and the address are the same, then nothing -- not even
			// a join operation -- is needed.
			if srv.Address == raft.ServerAddress(addr) && srv.ID == raft.ServerID(nodeId) {
				golog.Errorf("node %d at %s already member of cluster, ignoring join request", nodeId, addr)
				return nil
			}

			future := store.raft.RemoveServer(srv.ID, 0, 0)
			if err := future.Error(); err != nil {
				return fmt.Errorf("error removing existing node %d at %s: %s", nodeId, addr, err)
			}
		}
	}
	f := store.raft.AddVoter(raft.ServerID(nodeId), raft.ServerAddress(addr), 0, 0)
	if f.Error() != nil {
		return f.Error()
	}
	golog.Infof("node %d at %s joined successfully", nodeId, addr)
	return nil
}

func (store *Store) NodeID() uint64 {
	return store.nodeId
}

func (store *Store) ListenAddr() string {
	return store.listen
}

func (store *Store) IsLeader() bool {
	return store.raft.State() == raft.Leader
}

func (store *Store) Close() {
	snap := store.raft.Snapshot()
	if snap.Error() != nil && snap.Error().Error() != "nothing new to snapshot" {
		golog.Error(snap.Error())
	}
	store.raft.Shutdown()
	store.badger.Close()
	store.sqlstore.Close()
	store.server.Stop()
}
