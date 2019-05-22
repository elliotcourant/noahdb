package rpcwire

import (
	"fmt"
	"github.com/elliotcourant/noahdb/pkg/core"
	"github.com/elliotcourant/noahdb/pkg/pgproto"
	"github.com/readystock/golog"
	"io"
	"net"
)

type rpcWire struct {
	colony  core.Colony
	backend *pgproto.RpcBackend
}

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
				if err != io.EOF {
					golog.Errorf("failed serving rpc connection from [%s]: %v", conn.RemoteAddr(), err)
				}
			}
		}(colony, conn)
	}
}

func serveRpcConnection(colony core.Colony, conn net.Conn) error {
	backend, err := pgproto.NewRpcBackend(conn, conn)
	if err != nil {
		return err
	}

	// Right away let the RPC client know that we are ready for rpc calls.
	if err := backend.Send(&pgproto.ReadyForQuery{}); err != nil {
		return err
	}

	wire := &rpcWire{
		backend: backend,
		colony:  colony,
	}

	for {
		msg, err := backend.Receive()

		if err != nil {
			return err
		}

		switch message := msg.(type) {
		case *pgproto.JoinRequest:
			if err := wire.handleJoin(message); err != nil {
				backend.Send(&pgproto.ErrorResponse{Message: err.Error()})
			}
		case *pgproto.Terminate:
			return nil
		default:
			return fmt.Errorf("cannot handle message type: %v", message)
		}
	}
}
