package core

import (
	"fmt"
	"github.com/elliotcourant/noahdb/pkg/drivers/rpcer"
	"github.com/elliotcourant/noahdb/pkg/frunk"
	"github.com/elliotcourant/noahdb/pkg/transport"
	"github.com/elliotcourant/timber"
	"github.com/hashicorp/raft"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

type ColonyConfig struct {
	DataDirectory         string
	JoinAddresses         []raft.Server
	Transport             transport.TransportWrapper
	LocalPostgresAddress  string
	LocalPostgresPort     int32
	LocalPostgresUser     string
	LocalPostgresPassword string
	StartPool             bool
	AutoJoin              bool
}

type Accessors interface {
	Shards() ShardContext
	Tenants() TenantContext
	DataNodes() DataNodeContext
	Tables() TableContext
	Schema() SchemaContext
	Setting() SettingContext
	Users() UserContext
	Pool() PoolContext
	Sequences() SequenceContext
}

// Colony is a wrapper for all of the core data that noahdb needs to operate.
type Colony interface {
	Accessors

	Query(string) (*frunk.QueryResponse, error)

	CoordinatorID() uint64
	IsLeader() bool
	Close()
	Addr() net.Addr
	Join(id, addr string) error
	JoinCluster() error
	LeaderID() (string, string, error)
	State() frunk.ClusterState

	Neighbors() ([]*frunk.Server, error)

	InitColony(config ColonyConfig, log timber.Logger) error
}

func NewColony() Colony {
	return &base{}
}

func (ctx *base) InitColony(config ColonyConfig, log timber.Logger) error {
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	id := fmt.Sprintf("%s:%d", hostname, config.Transport.Port())

	fr := frunk.New(config.Transport.RaftTransport(), &frunk.StoreConfig{
		DBConf: &frunk.DBConfig{
			DSN:    "",
			Memory: false,
		},
		Dir: config.DataDirectory,
		ID:  id,
	})

	var potentialNeighbors []raft.Server
	if config.AutoJoin {
		potentialNeighbors, err = getAutoJoinAddresses()
		if err != nil {
			return err
		}
	}

	foundLeader, leaderAddress := false, ""

	// If we find other nodes in the cluster already, then we want to reach out to those nodes
	// and see if any of them are established. If we find one that is, send a join request to the
	// leader.
	for _, neighbor := range potentialNeighbors {
		rpcConn, err := rpcer.NewRPCDriver(id, config.Transport.Addr(), string(neighbor.Address))
		if err != nil {
			log.Errorf("failed to connect to potential neighbor [%s] at address %s: %v", neighbor.ID, neighbor.Address, err)
			continue
		}

		leaderAddr, err := rpcConn.Discover()
		if err != nil {
			log.Errorf("failed to discover via neighbor [%s] at address %s: %v", neighbor.ID, neighbor.Address, err)
			continue
		}

		if leaderAddr == "" {
			log.Warningf("neighbor [%s] at address [%s] is not established and has no leader", neighbor.ID, neighbor.Address)
			continue
		}

		foundLeader = true
		leaderAddress = leaderAddr
	}

	if leaderAddress != "" {

	}

	// joinAllowed, err := frunk.JoinAllowed(dataDirectory)
	// if err != nil {
	// 	return err
	// }

	// var joins []string
	// if joinAllowed {
	// 	joins, err = determineJoinAddresses(joinAddresses)
	// 	if err != nil {
	// 		return err
	// 	}
	// } else {
	// 	return err
	// }

	// Now, open store.
	if err := fr.Open(!foundLeader, potentialNeighbors...); err != nil {
		log.Fatalf("failed to open store: %s", err.Error())
	}

	*ctx = base{
		db:       fr,
		trans:    config.Transport,
		poolSync: sync.RWMutex{},
		pool:     map[uint64]*poolItem{},
	}

	if config.StartPool {
		ctx.Pool().StartPool()
	}

	if len(potentialNeighbors) > 0 {
		config.JoinAddresses = potentialNeighbors
	}

	if len(config.JoinAddresses) > 0 {
		attempts := 1
	RetryJoin:
		for i, joinAddr := range config.JoinAddresses {
			log.Debugf("trying to join address [%d] [%s]", i+1, joinAddr.Address)
			rpcDriver, err := rpcer.NewRPCDriver(id, config.Transport.Addr(), string(joinAddr.Address))
			if err != nil {
				log.Warningf("could not connect to join address [%s]: %v", joinAddr.Address, err)
				continue
			}
			if rpcDriver == nil {
				log.Warningf("failed to create frontend for address [%s]", joinAddr.Address)
				continue
			}
			if err := rpcDriver.Join(); err != nil {
				log.Warningf("could not join address [%s]: %v", joinAddr.Address, err)
				continue
			} else {
				log.Infof("successfully joined at address [%s]", joinAddr)
				goto WaitForSetup
			}
		}

		if attempts < 3 {
			log.Infof("was not able to join any of the nodes provided, will try again in 10 seconds; attempt: %d", attempts)
			time.Sleep(10 * time.Second)
			attempts++
			goto RetryJoin
		} else {
			log.Fatalf("failed to join any of the node found after %d attempt(s)", attempts)
		}
	}

WaitForSetup:

	openTimeout, err := time.ParseDuration("10s")
	if err != nil {
		log.Fatalf("failed to parse Raft open timeout: %s", err.Error())
	}
	fr.WaitForLeader(openTimeout)
	fr.WaitForApplied(openTimeout)

	// meta := map[string]string{}

	// // This may be a standalone server. In that case set its own metadata.
	// if err := fr.SetMetadata(meta); err != nil && err != store.ErrNotLeader {
	// 	// Non-leader errors are OK, since metadata will then be set through
	// 	// consensus as a result of a join. All other errors indicate a problem.
	// 	timber.Fatalf("failed to set store metadata: %s", err.Error())
	// }

	// time.Sleep(6 * time.Second)

	// handle joins here

	ctx.Setup(config)

	autoLocalPostgres(ctx, config)

	return nil
}

func autoLocalPostgres(ctx *base, config ColonyConfig) {
	if len(config.LocalPostgresAddress) > 0 && len(config.LocalPostgresUser) > 0 {
		// Check to see if there is a local postgres instance we can use.
		if _, err := ctx.DataNodes().NewDataNode(
			config.LocalPostgresAddress,
			config.LocalPostgresPort,
			config.LocalPostgresUser,
			config.LocalPostgresPassword); err != nil {
			timber.Warningf("could not add local postgres instance: %v", err)
			return
		}

		initialShards := 3
		for i := 0; i < initialShards; i++ {
			if _, err := ctx.Shards().NewShard(); err != nil {
				panic(err)
			}
		}

		if err := ctx.Shards().BalanceOrphanShards(); err != nil {
			panic(err)
		}
	}
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
