package testutils

import (
	"github.com/elliotcourant/noahdb/pkg/core"
	"github.com/elliotcourant/noahdb/pkg/pgwire"
	"github.com/elliotcourant/noahdb/pkg/rpcwire"
	"github.com/elliotcourant/noahdb/pkg/tcp"
	"github.com/readystock/golog"
	"io/ioutil"
	"net"
	"os"
	"strings"
)

func NewTestColonyEx(listenAddr string, joinAddresses ...string) (core.Colony, func()) {
	tempdir, err := ioutil.TempDir("", "core-temp")
	if err != nil {
		panic(err)
	}

	joins := strings.Join(joinAddresses, ",")

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

	go func() {
		if err = pgwire.NewServer(colony, trans); err != nil {
			golog.Errorf(err.Error())
		}
	}()

	go func() {
		if err = rpcwire.NewRpcServer(colony, trans); err != nil {
			golog.Errorf(err.Error())
		}
	}()

	err = colony.InitColony(tempdir, joins, trans)
	if err != nil {
		panic(err)
	}

	if colony.IsLeader() {
		_, err = colony.DataNodes().NewDataNode("127.0.0.1", os.Getenv("PGPASSWORD"), os.Getenv("PGPORT"))
		if err != nil {
			panic(err)
		}

		_, err = colony.Shards().NewShard()
		if err != nil {
			panic(err)
		}
	}

	return colony, func() {
		if err := os.RemoveAll(tempdir); err != nil {
			panic(err)
		}
	}
}

func NewTestColony(joinAddresses ...string) (core.Colony, func()) {
	return NewTestColonyEx(":", joinAddresses...)
}
