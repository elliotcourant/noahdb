package top

import (
	"fmt"
	"github.com/elliotcourant/noahdb/pkg/core"
	"github.com/elliotcourant/noahdb/pkg/pgwire"
	"github.com/readystock/golog"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type pgwireConfig struct {
	address string
	port    int
}

func (conf pgwireConfig) Address() string {
	return conf.address
}

func (conf pgwireConfig) Port() int {
	return conf.port
}

func NoahMain(dataDirectory, listenAddress, joinAddress, postgresAddress string) {
	golog.Debugf("starting noahdb")
	colony, err := core.NewColony(dataDirectory, listenAddress, joinAddress, postgresAddress)
	if err != nil {
		panic(fmt.Sprintf("could not setup colony: %s", err.Error()))
	} else if colony == nil {
		panic("failed to create a valid colony")
	}

	golog.Debugf("colony initialized, coordinator [%d]", colony.CoordinatorID())

	tasks := new(sync.WaitGroup)
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	signal.Notify(ch, os.Interrupt, syscall.SIGSEGV)
	signal.Notify(ch, os.Interrupt, syscall.SIGQUIT)

	tasks.Add(2)

	go func() {
		defer tasks.Done()
		<-ch
		golog.Warnf("stopping coordinator[%d]", colony.CoordinatorID())
		colony.Close()
		os.Exit(0)
	}()

	go func() {
		defer tasks.Done()
		if err = pgwire.NewServer(colony, pgwireConfig{
			address: "127.0.0.1",
			port:    5433,
		}); err != nil {
			golog.Errorf(err.Error())
		}
	}()

	tasks.Wait()
}
