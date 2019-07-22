package top

import (
	"fmt"
	"github.com/elliotcourant/noahdb/pkg/core"
	"github.com/elliotcourant/noahdb/pkg/kube"
	"github.com/elliotcourant/noahdb/pkg/pgwire"
	"github.com/elliotcourant/noahdb/pkg/rpcwire"
	"github.com/elliotcourant/noahdb/pkg/tcp"
	"github.com/elliotcourant/noahdb/pkg/util"
	"github.com/elliotcourant/timber"
	"github.com/hashicorp/raft"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

func NoahMain(dataDirectory, joinAddresses, listenAddr string, autoDataNode, autoJoin bool) {
	log := timber.New()

	log.Debugf("starting noahdb")
	l, err := util.ResolveAddress(listenAddr)
	if err != nil {
		panic(err)
	}
	listenAddr = l
	parsedRaftAddr, err := net.ResolveTCPAddr("tcp", listenAddr)
	if err != nil {
		panic(err)
		// return nil, nil, err
	}

	tn := tcp.NewTransport()

	if err := tn.Open(parsedRaftAddr.String()); err != nil {
		panic(err)
	}

	trans := core.NewTransportWrapper(tn)

	colony := core.NewColony()

	tasks := new(sync.WaitGroup)
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	signal.Notify(ch, os.Interrupt, syscall.SIGSEGV)
	signal.Notify(ch, os.Interrupt, syscall.SIGQUIT)

	tasks.Add(3)

	go func() {
		defer tasks.Done()
		<-ch
		log.Warningf("stopping coordinator[%d]", colony.CoordinatorID())
		colony.Close()
		os.Exit(0)
	}()

	go func() {
		defer tasks.Done()
		if err = pgwire.NewServer(colony, trans); err != nil {
			log.Errorf(err.Error())
		}
	}()

	go func() {
		defer tasks.Done()
		if err = rpcwire.NewRpcServer(colony, trans); err != nil {
			log.Errorf(err.Error())
		}
	}()

	go func() {
		defer tasks.Done()
		kube.RunEyeholes(colony)
	}()

	go func() {
		for {
			time.Sleep(30 * time.Second)
			if colony == nil {
				log.Infof("still strapping my boots")
			} else {
				addr, leaderId, err := colony.LeaderID()
				if err != nil {
					log.Errorf("could not get leader ID: %s", err)
				} else {
					log.Infof("current state [%s] current leader: %s | %s", colony.State(), leaderId, addr)
				}
			}
		}
	}()

	joins := make([]raft.Server, 0)

	// // If we are auto joining and no join address have been specified
	// if autoJoin && joinAddresses == "" {
	// 	potentialJoins, err := core.getAutoJoinAddresses()
	// 	if err != nil {
	// 		// If something went wrong inside the auto join address function, we likely
	// 		// would not be able to continue
	// 		panic(err)
	// 	}
	// 	golog.Infof("found %d potential auto-join addresses", len(potentialJoins))
	// 	joins = potentialJoins
	// } else {
	//
	// }

	if joinAddresses != "" {
		addresses := strings.Split(joinAddresses, ",")
		for _, addr := range addresses {
			if parsedAddress, err := util.ResolveAddress(addr); err != nil {
				log.Errorf("could not parse join address [%s]: %v", addr, err)
				panic(err)
			} else {
				joins = append(joins, raft.Server{
					ID:       raft.ServerID(parsedAddress),
					Address:  raft.ServerAddress(parsedAddress),
					Suffrage: raft.Voter,
				})
			}
		}
	}

	config := core.ColonyConfig{
		DataDirectory: dataDirectory,
		JoinAddresses: joins,
		Transport:     trans,
		AutoJoin:      autoJoin,
	}

	switch autoDataNode {
	case true:
		strPgPort, pgUser, pgPassword :=
			os.Getenv("PGPORT"), os.Getenv("PGUSER"), os.Getenv("PGPASSWORD")

		if strPgPort == "" || pgUser == "" {
			log.Infof("no auto data node could be found")
			break
		}

		addr, err := util.ExternalIP()
		if err != nil {
			log.Warningf("could not get external IP address for local data node")
			break
		}

		pgPort, err := strconv.ParseInt(strPgPort, 10, 32)
		if err != nil {
			log.Warningf("could not parse PGPORT environment variable: %s", strPgPort)
			break
		}

		config.LocalPostgresPort = int32(pgPort)
		config.LocalPostgresAddress = addr
		config.LocalPostgresPassword = pgPassword
		config.LocalPostgresUser = pgUser
	case false:
		timber.Verbosef("not using auto data node")
	}

	err = colony.InitColony(config, log)
	if err != nil {
		panic(fmt.Sprintf("could not setup colony: %s", err.Error()))
	} else if colony == nil {
		panic("failed to create a valid colony")
	}

	if colony.IsLeader() && autoDataNode {
		log.Infof("auto-detecting a local PostgreSQL instance")
	}

	log.Debugf("colony initialized, coordinator [%d]", colony.CoordinatorID())

	tasks.Wait()
}
