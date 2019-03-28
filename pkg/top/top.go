package top

import (
	"fmt"
	"github.com/elliotcourant/noahdb/pkg/core"
	"github.com/readystock/golog"
	"os"
	"os/signal"
	"syscall"
)

func NoahMain(dataDirectory, listenAddress, joinAddress, postgresAddress string) {
	golog.Debugf("starting noahdb")
	colony, err := core.NewColony(dataDirectory, listenAddress, joinAddress, postgresAddress)
	if err != nil {
		panic(fmt.Sprintf("could not setup colony: %s", err.Error()))
	} else if colony == nil {
		panic("failed to create a valid colony")
	}

	golog.Debugf("colony initialized, coordinator [%d]", colony.CoordinatorID())

	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	signal.Notify(ch, os.Interrupt, syscall.SIGSEGV)
	signal.Notify(ch, os.Interrupt, syscall.SIGQUIT)

	go func() {
		<-ch
		golog.Warnf("stopping coordinator[%d]", colony.CoordinatorID())
		colony.Close()
		os.Exit(0)
	}()
}
