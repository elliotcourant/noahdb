package pgtransport

import (
	"context"
	"github.com/elliotcourant/noahdb/pkg/pgproto"
	"net"
	"sync"
	"time"

	"github.com/elliotcourant/noahdb/pkg/logger"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/raft"
)

// deferError can be embedded to allow a future
// to provide an error in the future.
type deferError struct {
	err       error
	errCh     chan error
	responded bool
}

func (d *deferError) init() {
	d.errCh = make(chan error, 1)
}

func (d *deferError) Error() error {
	if d.err != nil {
		// Note that when we've received a nil error, this
		// won't trigger, but the channel is closed after
		// send so we'll still return nil below.
		return d.err
	}
	if d.errCh == nil {
		panic("waiting for response on nil channel")
	}
	d.err = <-d.errCh
	return d.err
}

func (d *deferError) respond(err error) {
	if d.errCh == nil {
		return
	}
	if d.responded {
		return
	}
	d.errCh <- err
	close(d.errCh)
	d.responded = true
}

// ServerAddressProvider just provides us a potential implementation
// to allow us to lookup an address with whatever ID we are provided.
// While it is default behavior most of the time to use the listen
// address as the server ID in a raft implementation, this is a dumb
// idea and we should absolutely not depend on it.
type ServerAddressProvider interface {
	ServerAddr(id raft.ServerID) (raft.ServerAddress, error)
}

// StreamLayer is just a local interface definition for our net stuff
// essentially what will actually be passed here is from the core.Wrapper
// stuff that we built as a net code hack.
type StreamLayer interface {
	net.Listener

	// Dial is used to create a new outgoing connection
	Dial(address raft.ServerAddress, timeout time.Duration) (net.Conn, error)
}

// PgTransport is an improved TCP transport for
// raft that uses a net code similar to Postgres.
type PgTransport struct {
	connPool     map[raft.ServerAddress][]*pgConn
	connPoolLock sync.Mutex

	consumeChannel chan raft.RPC

	hearbeatCallback      func(raft.RPC)
	heartbeatCallbackLock sync.Mutex

	// In the other TCP transport we use a different
	// logger, but to be consistent with what hashicorp's
	// raft library uses, we should use this.
	logger hclog.Logger

	maxPool int

	serverAddressProvider ServerAddressProvider

	shutdown        bool
	shutdownChannel chan struct{}
	shutdownLock    sync.Mutex

	stream StreamLayer

	streamContext     context.Context
	streamCancel      context.CancelFunc
	streamContextLock sync.RWMutex

	timeout      time.Duration
	timeoutScale int
}

// PgTransportConfig exposes just a few ways to tweak the
// internal behavior of the pg transport.
type PgTransportConfig struct {
	ServerAddressProvider ServerAddressProvider
	Logger                hclog.Logger
	Stream                StreamLayer
	MaxPool               int
	Timeout               time.Duration
}

// appendFuture is used for waiting on a pipelined append
// entries RPC.
type appendFuture struct {
	deferError
	start time.Time
	args  *raft.AppendEntriesRequest
	resp  *raft.AppendEntriesResponse
}

func (a *appendFuture) Start() time.Time {
	return a.start
}

func (a *appendFuture) Request() *raft.AppendEntriesRequest {
	return a.args
}

func (a *appendFuture) Response() *raft.AppendEntriesResponse {
	return a.resp
}

type pgConn struct {
	target raft.ServerAddress
	conn   net.Conn
}

func (p *pgConn) Release() error {
	return p.conn.Close()
}

type pgPipeline struct {
	conn      *pgConn
	transport *PgTransport

	doneChannel       chan raft.AppendFuture
	inProgressChannel chan appendFuture

	shutdown        bool
	shutdownChannel chan struct{}
	shutdownLock    sync.Mutex
}

func NewPgTransportWithConfig(
	config *PgTransportConfig,
) *PgTransport {
	if config.Logger == nil {
		config.Logger = logger.NewLogger()
	}
	trans := &PgTransport{
		connPool:              make(map[raft.ServerAddress][]*pgConn),
		consumeChannel:        make(chan raft.RPC),
		logger:                config.Logger,
		maxPool:               config.MaxPool,
		shutdownChannel:       make(chan struct{}),
		stream:                config.Stream,
		timeout:               config.Timeout, // I'm leary of this at the moment
		serverAddressProvider: config.ServerAddressProvider,
	}

	trans.setupStreamContext()
	go trans.listen()
	return trans
}

func (p *PgTransport) setupStreamContext() {
	p.streamContext, p.streamCancel = context.WithCancel(context.Background())
}

// getStreamContext is used retrieve the current stream context.
func (p *PgTransport) getStreamContext() context.Context {
	p.streamContextLock.RLock()
	defer p.streamContextLock.RUnlock()
	return p.streamContext
}

// LocalAddr implements the Transport interface.
func (p *PgTransport) LocalAddr() raft.ServerAddress {
	return raft.ServerAddress(p.stream.Addr().String())
}

func (p *PgTransport) IsShutdown() bool {
	select {
	case <-p.shutdownChannel:
		return true
	default:
		return false
	}
}

func (p *PgTransport) listen() {
	for {
		conn, err := p.stream.Accept()
		if err != nil {
			if p.IsShutdown() {
				return
			}
			p.logger.Error("failed to accept connection: %v", err)
			continue
		}
		p.logger.Trace("%v accepted connection from: %v", p.LocalAddr(), conn.RemoteAddr())
	}
}

func (p *PgTransport) handleConnection(connectionContext context.Context, conn net.Conn) {
	defer conn.Close()
	wire, err := pgproto.NewRaftWire(conn, conn)
	if err != nil {
		p.logger.Error("could not create raft wire for connection [%v]: %v", conn.RemoteAddr(), err)
	}

	for {
		select {
		case <-connectionContext.Done():
			p.logger.Debug("stream layer is closed")
			return
		default:
		}

		request, err := wire.Receive()
		if err != nil {
			p.logger.Error("failed to receive message from [%v]: %v", conn.RemoteAddr(), err)
			return
		}

		responseChannel := make(chan raft.RPCResponse, 1)
		rpc := raft.RPC{
			RespChan: responseChannel,
		}

		isHeartbeat := false
		switch req := request.(type) {
		case *pgproto.AppendEntriesRequest:
			rpc.Command = req.AppendEntriesRequest

			// Check if this is a heartbeat
			if req.Term != 0 && req.Leader != nil &&
				req.PrevLogEntry == 0 && req.PrevLogTerm == 0 &&
				len(req.Entries) == 0 && req.LeaderCommitIndex == 0 {
				isHeartbeat = true
			}
		case *pgproto.RequestVoteRequest:
			rpc.Command = req.RequestVoteRequest
		case *pgproto.InstallSnapshotRequest:
			rpc.Command = req.InstallSnapshotRequest
			rpc.Reader = req.Reader()
		default:
			p.logger.Error("did not recognize request type [%v] from [%v]: %v", req, conn.RemoteAddr(), err)
			return
		}

		// Check for heartbeat fast-path
		if isHeartbeat {
			p.heartbeatCallbackLock.Lock()
			callback := p.hearbeatCallback
			p.heartbeatCallbackLock.Unlock()
			if callback != nil {
				callback(rpc)
				goto RESPONSE
			}
		}

		// Dispatch the RPC to this raft node.
		select {
		case p.consumeChannel <- rpc:
		case <-p.shutdownChannel:
			p.logger.Error("transport is shutdown")
			return
		}

		// Wait for response
	RESPONSE:
		select {
		case response := <-responseChannel:
			var msg pgproto.Message
			switch rsp := response.Response.(type) {
			case raft.AppendEntriesResponse:
				msg = &pgproto.AppendEntriesResponse{
					Error:                 response.Error,
					AppendEntriesResponse: rsp,
				}
			case raft.RequestVoteResponse:
				msg = &pgproto.RequestVoteResponse{
					Error:               response.Error,
					RequestVoteResponse: rsp,
				}
			case raft.InstallSnapshotResponse:
				msg = &pgproto.InstallSnapshotResponse{
					Error:                   response.Error,
					InstallSnapshotResponse: rsp,
				}
			case nil:
				msg = &pgproto.ErrorResponse{
					Message: response.Error.Error(),
				}
			}

			if _, err := conn.Write(msg.Encode(nil)); err != nil {
				p.logger.Error("failed to send response to [%v]: %v", conn.RemoteAddr(), err)
				return
			}
		case <-p.shutdownChannel:
			p.logger.Warn("closing transport due to shutdown")
			return
		}
	}
}
