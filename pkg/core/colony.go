package core

import (
	"fmt"
	"github.com/elliotcourant/noahdb/pkg/drivers/rpcer"
	"github.com/elliotcourant/noahdb/pkg/frunk"
	"github.com/readystock/golog"
	"net"
	"os"
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
	IsLeader() bool
	Close()
	Addr() net.Addr
	Join(id, addr string) error
	JoinCluster() error
	LeaderID() (string, error)
	State() frunk.ClusterState

	Neighbors() ([]*frunk.Server, error)

	InitColony(dataDirectory, joinAddresses string, trans TransportWrapper) error
}

func NewColony() Colony {
	return &base{}
}

func (ctx *base) InitColony(dataDirectory, joinAddresses string, trans TransportWrapper) error {
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	id := fmt.Sprintf("%s:%d", hostname, trans.Port())

	fr := frunk.New(trans.RaftTransport(), &frunk.StoreConfig{
		DBConf: &frunk.DBConfig{
			DSN:    "",
			Memory: true,
		},
		Dir: dataDirectory,
		ID:  id,
	})

	joinAllowed, err := frunk.JoinAllowed(dataDirectory)
	if err != nil {
		return err
	}

	var joins []string
	if joinAllowed {
		joins, err = determineJoinAddresses(joinAddresses)
		if err != nil {
			return err
		}
	} else {
		return err
	}

	// Now, open store.
	if err := fr.Open(len(joins) == 0); err != nil {
		golog.Fatalf("failed to open store: %s", err.Error())
	}

	*ctx = base{
		db:       fr,
		trans:    trans,
		poolSync: sync.Mutex{},
		pool:     map[uint64]*poolItem{},
	}

	openTimeout, err := time.ParseDuration("10s")
	if err != nil {
		golog.Fatalf("failed to parse Raft open timeout: %s", err.Error())
	}
	fr.WaitForLeader(openTimeout)
	fr.WaitForApplied(openTimeout)

	// meta := map[string]string{}

	// // This may be a standalone server. In that case set its own metadata.
	// if err := fr.SetMetadata(meta); err != nil && err != store.ErrNotLeader {
	// 	// Non-leader errors are OK, since metadata will then be set through
	// 	// consensus as a result of a join. All other errors indicate a problem.
	// 	golog.Fatalf("failed to set store metadata: %s", err.Error())
	// }

	// time.Sleep(6 * time.Second)

	// handle joins here

	if len(joins) > 0 {
		for i, joinAddr := range joins {
			golog.Debugf("trying to join address [%d] [%s]", i+1, joinAddr)
			rpcDriver, err := rpcer.NewRPCDriver(id, trans.Addr(), joinAddr)
			if err != nil {
				golog.Warnf("could not connect to join address [%s]: %v", joinAddr, err)
			}
			if rpcDriver == nil {
				golog.Warnf("failed to create frontend for address [%s]", joinAddr)
			}
			if err := rpcDriver.Join(); err != nil {
				golog.Warnf("could not join address [%s]: %v", joinAddr, err)
			} else {
				golog.Infof("successfully joined at address [%s]", joinAddr)
				break
			}
		}
	}

	ctx.Setup()

	return nil
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
