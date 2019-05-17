package top

import (
	"fmt"
	"github.com/elliotcourant/noahdb/pkg/core"
	"github.com/elliotcourant/noahdb/pkg/kube"
	"github.com/elliotcourant/noahdb/pkg/pgwire"
	"github.com/elliotcourant/noahdb/pkg/rpcwire"
	"github.com/elliotcourant/noahdb/pkg/tcp"
	"github.com/readystock/golog"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func NoahMain(dataDirectory, joinAddresses, listenAddr string, autoDataNode bool) {
	golog.Debugf("starting noahdb")

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

	err = colony.InitColony(dataDirectory, joinAddresses, trans)
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
