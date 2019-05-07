package transport

import (
	"bufio"
	"context"
	"github.com/hashicorp/raft"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

/*

PgTransport provides a network based transport that can be
used to communicate with Raft on remote machines. It requires
an underlying stream layer to provide a stream abstraction, which can
be simple TCP, TLS, etc.

This transport is very simple and lightweight. Each RPC request is
framed by sending a byte that indicates the message type, followed
by the MsgPack encoded request.

The response is an error string followed by the response object,
both are encoded using MsgPack.

InstallSnapshot is special, in that after the RPC request we stream
the entire state. That socket is not re-used as the connection state
is not known if there is an error.

*/
type PgTransport struct {
	connPool     map[raft.ServerAddress][]*frontendPgConn
	connPoolLock sync.Mutex

	consumeCh chan raft.RPC

	heartbeatFn     func(raft.RPC)
	heartbeatFnLock sync.Mutex

	logger *log.Logger

	maxPool int

	serverAddressProvider ServerAddressProvider

	shutdown     bool
	shutdownCh   chan struct{}
	shutdownLock sync.Mutex

	stream StreamLayer

	// streamCtx is used to cancel existing connection handlers.
	streamCtx     context.Context
	streamCancel  context.CancelFunc
	streamCtxLock sync.RWMutex

	timeout      time.Duration
	TimeoutScale int
}

// PgTransportConfig encapsulates configuration for the network transport layer.
type PgTransportConfig struct {
	// ServerAddressProvider is used to override the target address when establishing a connection to invoke an RPC
	ServerAddressProvider ServerAddressProvider

	Logger *log.Logger

	// Dialer
	Stream StreamLayer

	// MaxPool controls how many connections we will pool
	MaxPool int

	// Timeout is used to apply I/O deadlines. For InstallSnapshot, we multiply
	// the timeout by (SnapshotSize / TimeoutScale).
	Timeout time.Duration
}

type frontendPgConn struct {
	target raft.ServerAddress
	conn   net.Conn
	r      *bufio.Reader
	w      *bufio.Writer
}

func (n *frontendPgConn) Release() error {
	return n.conn.Close()
}

// NewNetworkTransportWithConfig creates a new network transport with the given config struct
func NewPgTransportWithConfig(
	config *PgTransportConfig,
) *PgTransport {
	if config.Logger == nil {
		config.Logger = log.New(os.Stderr, "", log.LstdFlags)
	}
	trans := &PgTransport{
		connPool:              make(map[raft.ServerAddress][]*frontendPgConn),
		consumeCh:             make(chan raft.RPC),
		logger:                config.Logger,
		maxPool:               config.MaxPool,
		shutdownCh:            make(chan struct{}),
		stream:                config.Stream,
		timeout:               config.Timeout,
		TimeoutScale:          DefaultTimeoutScale,
		serverAddressProvider: config.ServerAddressProvider,
	}

	// Create the connection context and then start our listener.
	trans.setupStreamContext()
	go trans.listen()

	return trans
}

// setupStreamContext is used to create a new stream context. This should be
// called with the stream lock held.
func (n *PgTransport) setupStreamContext() {
	ctx, cancel := context.WithCancel(context.Background())
	n.streamCtx = ctx
	n.streamCancel = cancel
}
