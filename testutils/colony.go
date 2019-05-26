package testutils

import (
	"fmt"
	"github.com/elliotcourant/noahdb/pkg/core"
	"github.com/elliotcourant/noahdb/pkg/pgwire"
	"github.com/elliotcourant/noahdb/pkg/rpcwire"
	"github.com/elliotcourant/noahdb/pkg/tcp"
	"github.com/hashicorp/raft"
	"github.com/readystock/golog"
	"io/ioutil"
	"net"
	"os"
)

func NewTestColonyEx(listenAddr string, joinAddresses ...string) (core.Colony, func()) {
	golog.SetLevel("trace")
	tempdir, err := ioutil.TempDir("", "core-temp")
	if err != nil {
		panic(err)
	}

	// joins := strings.Join(joinAddresses, ",")

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

	err = colony.InitColony(tempdir, make([]raft.Server, 0), trans)
	if err != nil {
		panic(err)
	}

	// if colony.IsLeader() {
	// 	_, err = colony.DataNodes().NewDataNode("127.0.0.1", os.Getenv("PGPASSWORD"), os.Getenv("PGPORT"))
	// 	if err != nil {
	// 		panic(err)
	// 	}
	//
	// 	_, err = colony.Shards().NewShard()
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }

	return colony, func() {
		if err := os.RemoveAll(tempdir); err != nil {
			panic(err)
		}
	}
}

func NewTestColony(joinAddresses ...string) (core.Colony, func()) {
	return NewTestColonyEx(":", joinAddresses...)
}

func ConnectionString(address net.Addr) string {
	addr, err := net.ResolveTCPAddr(address.Network(), address.String())
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		addr.IP.String(), addr.Port, "noah", "password", "postgres")
}
