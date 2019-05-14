package rpcwire

import (
	"github.com/elliotcourant/noahdb/pkg/core"
	"github.com/readystock/golog"
)

func NewRpcServer(colony core.Colony, transport core.TransportWrapper) error {
	ln := transport.RpcTransport()

	for {
		conn, err := ln.Accept()
		if err != nil {
			golog.Errorf("could not accept rpc connection: %v", err)
			continue
		}

		golog.Verbosef("received rpc connection from [%s]", conn.RemoteAddr().String())

		conn.Close()
	}
}
