package core

import (
	"github.com/readystock/golog"
	"net"
	"time"
)

// Listener is the interface Raft-compatible network layers
// should implement.
type Listener interface {
	net.Listener
	Dial(address string, timeout time.Duration) (net.Conn, error)
}

type TransportWrapper interface {
	NormalTransport() net.Listener
	ForwardToRaft(net.Conn, error)
	ForwardToRpc(net.Conn, error)
	RaftTransport() Listener
	RpcTransport() Listener
	Addr() net.Addr
	Close()
}

type accept struct {
	conn net.Conn
	err  error
}
type transportWrapperItem struct {
	listener      Listener
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
	transport     Listener
	raftTransport *transportWrapperItem
	rpcTransport  *transportWrapperItem
}

func NewTransportWrapper(listener Listener) TransportWrapper {
	wrapper := &transportWrapper{
		transport: listener,
		raftTransport: &transportWrapperItem{
			acceptChannel: make(chan accept, 0),
		},
		rpcTransport: &transportWrapperItem{
			acceptChannel: make(chan accept, 0),
		},
	}

	{
		wrapper.raftTransport.closeCallback = wrapper.closeCallback
		wrapper.raftTransport.listener = wrapper.transport
	}

	{
		wrapper.rpcTransport.closeCallback = wrapper.closeCallback
		wrapper.rpcTransport.listener = wrapper.transport
	}

	return wrapper
}

func (wrapper *transportWrapper) ForwardToRaft(conn net.Conn, err error) {
	wrapper.raftTransport.SendAccept(conn, err)
}

func (wrapper *transportWrapper) ForwardToRpc(conn net.Conn, err error) {
	wrapper.rpcTransport.SendAccept(conn, err)
}

func (wrapper *transportWrapper) closeCallback() {
	golog.Verbosef("received close callback")
}

func (wrapper *transportWrapper) RaftTransport() Listener {
	return wrapper.raftTransport
}

func (wrapper *transportWrapper) RpcTransport() Listener {
	return wrapper.rpcTransport
}

func (wrapper *transportWrapper) NormalTransport() net.Listener {
	return wrapper.transport
}

func (wrapper *transportWrapper) Addr() net.Addr {
	return wrapper.transport.Addr()
}

func (wrapper *transportWrapper) Close() {
	wrapper.raftTransport.Close()
}
