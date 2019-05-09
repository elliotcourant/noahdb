package core

import (
	"github.com/elliotcourant/noahdb/pkg/frunk"
	"github.com/elliotcourant/noahdb/pkg/store"
	"github.com/elliotcourant/noahdb/pkg/tcp"
	"github.com/readystock/golog"
	"net"
	"strings"
	"sync"
	"time"
)

// Colony is a wrapper for all of the core data that noahdb needs to operate.
type Colony interface {
	Shards() ShardContext
	Tenants() TenantContext
	DataNodes() DataNodeContext
	Tables() TableContext
	Schema() SchemaContext
	Users() UserContext
	Pool() PoolContext
	Query(string) (*frunk.QueryResponse, error)
	// Shards()
	// Nodes()
	// Tenants()
	// Network()
	// Settings()
	// Pool()
	// Schema()
	// Sequences()

	CoordinatorID() uint64
	Close()
}

func NewColony(dataDirectory, joinAddresses, postgresAddress, raftAddr string) (Colony, error) {
	// db, err := store.CreateStore(dataDirectory, listenAddress, "")
	// if err != nil {
	// 	return nil, err
	// }

	parsedRaftAddr, err := net.ResolveTCPAddr("tcp", raftAddr)
	if err != nil {
		return nil, err
	}

	tn := tcp.NewTransport()

	if err := tn.Open(parsedRaftAddr.String()); err != nil {
		return nil, err
	}

	tn.Accept()

	fr := frunk.New(tn, &frunk.StoreConfig{
		DBConf: &frunk.DBConfig{
			DSN:    "",
			Memory: true,
		},
		Dir: dataDirectory,
		ID:  parsedRaftAddr.String(),
	})

	joinAllowed, err := frunk.JoinAllowed(dataDirectory)
	if err != nil {
		return nil, err
	}

	var joins []string
	if joinAllowed {
		joins, err = determineJoinAddresses(joinAddresses)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}

	// Now, open store.
	if err := fr.Open(len(joins) == 0); err != nil {
		golog.Fatalf("failed to open store: %s", err.Error())
	}

	// handle joins here

	openTimeout, err := time.ParseDuration("120s")
	if err != nil {
		golog.Fatalf("failed to parse Raft open timeout: %s", err.Error())
	}
	fr.WaitForLeader(openTimeout)
	fr.WaitForApplied(openTimeout)

	meta := map[string]string{}

	// This may be a standalone server. In that case set its own metadata.
	if err := fr.SetMetadata(meta); err != nil && err != store.ErrNotLeader {
		// Non-leader errors are OK, since metadata will then be set through
		// consensus as a result of a join. All other errors indicate a problem.
		golog.Fatalf("failed to set store metadata: %s", err.Error())
	}

	colony := &base{
		db:       fr,
		poolSync: sync.Mutex{},
		pool:     map[uint64]*poolItem{},
	}

	time.Sleep(6 * time.Second)

	colony.Setup()

	return colony, nil
}

func determineJoinAddresses(joinAddr string) ([]string, error) {

	var addrs []string
	if joinAddr != "" {
		// Explicit join addresses are first priority.
		addrs = strings.Split(joinAddr, ",")
	}

	// if discoID != "" {
	// 	log.Printf("registering with Discovery Service at %s with ID %s", discoURL, discoID)
	// 	c := disco.New(discoURL)
	// 	r, err := c.Register(discoID, apiAdv)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	log.Println("Discovery Service responded with nodes:", r.Nodes)
	// 	for _, a := range r.Nodes {
	// 		if a != apiAdv {
	// 			// Only other nodes can be joined.
	// 			addrs = append(addrs, a)
	// 		}
	// 	}
	// }

	return addrs, nil
}
