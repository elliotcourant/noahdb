package rpcwire

import (
	"fmt"
	"github.com/elliotcourant/noahdb/pkg/core"
	"github.com/elliotcourant/noahdb/pkg/pgproto"
	"github.com/readystock/golog"
	"net"
)

func NewRpcServer(colony core.Colony, transport core.TransportWrapper) error {
	ln := transport.RpcTransport()

	for {
		conn, err := ln.Accept()
		if err != nil {
			golog.Errorf("could not accept rpc connection: %v", err)
			continue
		}

		go func(colony core.Colony, conn net.Conn) {
			defer conn.Close()
			if err := serveRpcConnection(colony, conn); err != nil {
				golog.Errorf("failed serving rpc connection: %v", err)
			}
		}(colony, conn)
	}
}

func serveRpcConnection(colony core.Colony, conn net.Conn) error {
	golog.Verbosef("received rpc connection from [%s]", conn.RemoteAddr().String())

	backend, err := pgproto.NewBackend(conn, conn)
	if err != nil {
		return err
	}

	// Right away let the RPC client know that we are ready for rpc calls.
	if err := backend.Send(&pgproto.ReadyForQuery{}); err != nil {
		return err
	}

	msg, err := backend.Receive()

	if err != nil {
		return err
	}

	switch msg.(type) {
	case *pgproto.JoinRequest:

	default:
		return fmt.Errorf("cannot handle message type: %v", msg)
	}
	return nil
}
