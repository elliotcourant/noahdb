package top

import (
	"fmt"
	"github.com/elliotcourant/noahdb/pkg/core"
	"github.com/elliotcourant/noahdb/pkg/kube"
	"github.com/elliotcourant/noahdb/pkg/pgwire"
	"github.com/elliotcourant/noahdb/pkg/rpcwire"
	"github.com/elliotcourant/noahdb/pkg/tcp"
	"github.com/elliotcourant/noahdb/pkg/util"
	"github.com/hashicorp/raft"
	"github.com/readystock/golog"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

func NoahMain(dataDirectory, joinAddresses, listenAddr string, autoDataNode, autoJoin bool) {
	golog.Debugf("starting noahdb")
	l, err := util.ResolveLocalAddress(listenAddr)
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
		golog.Warnf("stopping coordinator[%d]", colony.CoordinatorID())
		colony.Close()
		os.Exit(0)
	}()

	go func() {
		defer tasks.Done()
		if err = pgwire.NewServer(colony, trans); err != nil {
			golog.Errorf(err.Error())
		}
	}()

	go func() {
		defer tasks.Done()
		if err = rpcwire.NewRpcServer(colony, trans); err != nil {
			golog.Errorf(err.Error())
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
				golog.Infof("still strapping my boots")
			} else {
				addr, leaderId, err := colony.LeaderID()
				if err != nil {
					golog.Errorf("could not get leader ID: %s", err)
				} else {
					golog.Infof("current state [%s] current leader: %s | %s", colony.State(), leaderId, addr)
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
			if parsedAddress, err := util.ResolveLocalAddress(addr); err != nil {
				golog.Errorf("could not parse join address [%s]: %v", addr, err)
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

	err = colony.InitColony(dataDirectory, joins, trans)
	if err != nil {
		panic(fmt.Sprintf("could not setup colony: %s", err.Error()))
	} else if colony == nil {
		panic("failed to create a valid colony")
	}

	if colony.IsLeader() && autoDataNode {
		golog.Infof("auto-detecting a local PostgreSQL instance")
	}

	golog.Debugf("colony initialized, coordinator [%d]", colony.CoordinatorID())

	tasks.Wait()
}
