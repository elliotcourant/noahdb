package core

import (
	"github.com/elliotcourant/noahdb/pkg/pgproto"
	"github.com/readystock/golog"
	"io"
	"net"
	"time"
)

// Listener is the interface Raft-compatible network layers
// should implement.
type Listener interface {
	net.Listener
	Dial(address string, timeout time.Duration) (net.Conn, error)
}

type TransportPredicate func(io.Reader) bool

type accept struct {
	conn net.Conn
	err  error
}
type transportWrapperItem struct {
	listener  Listener
	predicate TransportPredicate

	acceptChannel chan accept

	closeCallback func()
}

func (t *transportWrapperItem) SendAccept(conn net.Conn, err error) {
	t.acceptChannel <- accept{conn, err}
}

func (t *transportWrapperItem) Accept() (net.Conn, error) {
	a := <-t.acceptChannel
	return a.conn, a.err
}

func (t *transportWrapperItem) Close() error {
	close(t.acceptChannel)
	t.closeCallback()
	return t.listener.Close()
}

func (t *transportWrapperItem) Addr() net.Addr {
	return t.listener.Addr()
}

func (t *transportWrapperItem) Dial(address string, timeout time.Duration) (net.Conn, error) {
	return t.listener.Dial(address, timeout)
}

type transportWrapper struct {
	transport Listener

	raftTransport     *transportWrapperItem
	rpcTransport      *transportWrapperItem
	postgresTransport *transportWrapperItem
}

func (wrapper *transportWrapper) Start() {
	sendError := func(err error) {
		// If there is an error, issue the error to all the listening connections.
		wrapper.raftTransport.SendAccept(nil, err)
		wrapper.rpcTransport.SendAccept(nil, err)
		wrapper.postgresTransport.SendAccept(nil, err)
	}

	go func() {
		for {
			conn, err := wrapper.transport.Accept()
			if err != nil {
				sendError(err)
				golog.Warnf("received error when accepting connection")
				continue
			}

			backend, err := pgproto.NewBackend(conn, conn)
			if err != nil {
				sendError(err)
				golog.Warnf("received error when creating backend")
				continue
			}

			initial, err := backend.ReceiveInitialMessage()
			if err != nil {
				sendError(err)
				golog.Warnf("could not receive initial message")
			}

			switch initial {
			case pgproto.RaftNumber:
				wrapper.raftTransport.SendAccept(conn, nil)
			case pgproto.RpcNumber:
				wrapper.rpcTransport.SendAccept(conn, nil)
			case pgproto.ProtocolVersionNumber:

			default:

			}
		}
	}()
}

func (wrapper *transportWrapper) closeCallback() {
	golog.Verbosef("received close callback")
}
