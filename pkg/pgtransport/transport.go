package pgtransport

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/elliotcourant/noahdb/pkg/pgproto"
	"github.com/readystock/golog"
	"io"
	"net"
	"os"
	"sync"
	"time"

	"github.com/elliotcourant/noahdb/pkg/logger"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/raft"
)

const (
	// rpcMaxPipeline controls the maximum number of outstanding
	// AppendEntries RPC calls.
	rpcMaxPipeline = 128
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
	shutdownLock    sync.RWMutex

	stream StreamLayer

	streamContext     context.Context
	streamCancel      context.CancelFunc
	streamContextLock sync.RWMutex

	timeout      time.Duration
	timeoutScale int
}

func (p *PgTransport) TimeoutNow(id raft.ServerID, target raft.ServerAddress, args *raft.TimeoutNowRequest, resp *raft.TimeoutNowResponse) error {
	panic("implement me")
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
	wire   *pgproto.RaftWire
}

func (p *pgConn) Release() error {
	return p.conn.Close()
}

type pgPipeline struct {
	conn      *pgConn
	transport *PgTransport

	doneChannel       chan raft.AppendFuture
	inProgressChannel chan *appendFuture

	shutdown        bool
	shutdownChannel chan struct{}
	shutdownLock    sync.Mutex

	logger hclog.Logger

	wire *pgproto.RaftWire
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

// NewPgTransport creates a new network transport with the given dialer
// and listener. The maxPool controls how many connections we will pool. The
// timeout is used to apply I/O deadlines. For InstallSnapshot, we multiply
// the timeout by (SnapshotSize / TimeoutScale).
func NewPgTransport(
	stream StreamLayer,
	maxPool int,
	timeout time.Duration,
	logOutput io.Writer,
) *PgTransport {
	if logOutput == nil {
		logOutput = os.Stderr
	}
	config := &PgTransportConfig{Stream: stream, MaxPool: maxPool, Timeout: timeout, Logger: logger.NewLogger()}
	return NewPgTransportWithConfig(config)
}

// NewPgTransportWithLogger creates a new network transport with the given logger, dialer
// and listener. The maxPool controls how many connections we will pool. The
// timeout is used to apply I/O deadlines. For InstallSnapshot, we multiply
// the timeout by (SnapshotSize / TimeoutScale).
func NewPgTransportWithLogger(
	stream StreamLayer,
	maxPool int,
	timeout time.Duration,
	logger hclog.Logger,
) *PgTransport {
	config := &PgTransportConfig{Stream: stream, MaxPool: maxPool, Timeout: timeout, Logger: logger}
	return NewPgTransportWithConfig(config)
}

// SetHeartbeatHandler is used to setup a heartbeat handler
// as a fast-pass. This is to avoid head-of-line blocking from
// disk IO.
func (p *PgTransport) SetHeartbeatHandler(callback func(rpc raft.RPC)) {
	p.heartbeatCallbackLock.Lock()
	defer p.heartbeatCallbackLock.Unlock()
	p.hearbeatCallback = callback
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

// CloseStreams closes the current streams.
func (p *PgTransport) CloseStreams() {
	p.connPoolLock.Lock()
	defer p.connPoolLock.Unlock()

	for k, e := range p.connPool {
		for _, conn := range e {
			conn.Release()
		}

		delete(p.connPool, k)
	}

	// Cancel the existing connections and create a new context. Both these
	// operations must always be done with the lock held otherwise we can create
	// connection handlers that are holding a context that will never be
	// cancelable.
	p.streamContextLock.Lock()
	defer p.streamContextLock.Unlock()
	p.streamCancel()
	p.setupStreamContext()
}

func (p *PgTransport) Close() error {
	p.shutdownLock.Lock()
	defer p.shutdownLock.Unlock()

	if !p.shutdown {
		close(p.shutdownChannel)
		p.stream.Close()
		p.shutdown = true
	}
	return nil
}

// Consumer implements the Transport interface.
func (p *PgTransport) Consumer() <-chan raft.RPC {
	return p.consumeChannel
}

func (p *PgTransport) AppendEntriesPipeline(
	id raft.ServerID,
	target raft.ServerAddress,
) (raft.AppendPipeline, error) {
	conn, err := p.getConnFromAddressProvider(id, target)
	if err != nil {
		return nil, err
	}

	return newPgPipeline(p, conn)
}

func (p *PgTransport) AppendEntries(
	id raft.ServerID,
	target raft.ServerAddress,
	args *raft.AppendEntriesRequest,
	resp *raft.AppendEntriesResponse,
) error {
	return p.genericRPC(id, target, &pgproto.AppendEntriesRequest{
		AppendEntriesRequest: *args,
	}, resp)
}

func (p *PgTransport) RequestVote(
	id raft.ServerID,
	target raft.ServerAddress,
	args *raft.RequestVoteRequest,
	resp *raft.RequestVoteResponse,
) error {
	return p.genericRPC(id, target, &pgproto.RequestVoteRequest{
		RequestVoteRequest: *args,
	}, resp)
}

func (p *PgTransport) InstallSnapshot(
	id raft.ServerID,
	target raft.ServerAddress,
	args *raft.InstallSnapshotRequest,
	resp *raft.InstallSnapshotResponse,
	data io.Reader,
) error {
	snapshot := make([]byte, 0)
	writer := bytes.NewBuffer(snapshot)

	if _, err := io.Copy(writer, data); err != nil {
		p.logger.Error("failed to copy snapshot data to bytes: %v", err)
		return err
	}

	// Get a conn, always close for InstallSnapshot
	conn, err := p.getConnFromAddressProvider(id, target)
	if err != nil {
		return err
	}
	defer conn.Release()

	// Set a deadline, scaled by request size
	if p.timeout > 0 {
		timeout := p.timeout * time.Duration(args.Size/int64(p.timeoutScale))
		if timeout < p.timeout {
			timeout = p.timeout
		}
		conn.conn.SetDeadline(time.Now().Add(timeout))
	}

	if err := conn.wire.Send(&pgproto.InstallSnapshotRequest{
		InstallSnapshotRequest: *args,
		SnapshotData:           snapshot,
	}); err != nil {
		p.logger.Error("failed sending snapshot to [%v]: %v", conn.conn.RemoteAddr(), err)
		return err
	}

	return p.receiveResponse(conn, resp)
}

// EncodePeer implements the Transport interface.
func (p *PgTransport) EncodePeer(id raft.ServerID, a raft.ServerAddress) []byte {
	address := p.getProviderAddressOrFallback(id, a)
	return []byte(address)
}

// DecodePeer implements the Transport interface.
func (p *PgTransport) DecodePeer(buf []byte) raft.ServerAddress {
	return raft.ServerAddress(buf)
}

func (p *PgTransport) genericRPC(id raft.ServerID, target raft.ServerAddress, args pgproto.RaftMessage, response interface{}) error {
	conn, err := p.getConnFromAddressProvider(id, target)
	if err != nil {
		return err
	}

	if p.timeout > 0 {
		conn.conn.SetDeadline(time.Now().Add(p.timeout))
	}

	if err := conn.wire.Send(args); err != nil {
		p.logger.Error("when sending RPC to [%v]: %v", conn.conn.RemoteAddr(), err)
		return err
	}

	return p.receiveResponse(conn, response)
}

func (p *PgTransport) receiveResponse(conn *pgConn, response interface{}) error {
	message, err := conn.wire.Receive()
	if err != nil {
		if err == io.EOF || p.IsShutdown() {
			return nil
		}
		p.logger.Error("could not receive message from [%v]: %v", conn.conn.RemoteAddr(), err)
		return err
	}

	defer p.returnConn(conn)

	return func(message pgproto.RaftMessage, response interface{}) (err error) {
		switch msg := message.(type) {
		case *pgproto.AppendEntriesResponse:
			err = msg.Error
			if r, ok := response.(*raft.AppendEntriesResponse); ok {
				*r = msg.AppendEntriesResponse
			} else {
				p.logger.Warn("received %T but was expecting to received %T", msg, response)
			}
		case *pgproto.RequestVoteResponse:
			err = msg.Error
			if r, ok := response.(*raft.RequestVoteResponse); ok {
				*r = msg.RequestVoteResponse
			} else {
				p.logger.Warn("received %T but was expecting to received %T", msg, response)
			}
			response = msg.RequestVoteResponse
		case *pgproto.InstallSnapshotResponse:
			err = msg.Error
			if r, ok := response.(*raft.InstallSnapshotResponse); ok {
				*r = msg.InstallSnapshotResponse
			} else {
				p.logger.Warn("received %T but was expecting to received %T", msg, response)
			}
			response = msg.InstallSnapshotResponse
		case *pgproto.ErrorResponse:
			err = errors.New(msg.Message)
			response = nil
		default:
			err = fmt.Errorf("could handle response message type [%v]", msg)
			response = nil
		}
		return err
	}(message, response)
}

func (p *PgTransport) getPooledConn(target raft.ServerAddress) *pgConn {
	p.connPoolLock.Lock()
	defer p.connPoolLock.Unlock()

	conns, ok := p.connPool[target]
	if !ok || len(conns) == 0 {
		return nil
	}

	var conn *pgConn
	num := len(conns)
	conn, conns[num-1] = conns[num-1], nil
	p.connPool[target] = conns[:num-1]
	return conn
}

func (p *PgTransport) getConnFromAddressProvider(
	id raft.ServerID,
	target raft.ServerAddress,
) (*pgConn, error) {
	address := p.getProviderAddressOrFallback(id, target)
	return p.getConn(address)
}

func (p *PgTransport) getProviderAddressOrFallback(
	id raft.ServerID,
	target raft.ServerAddress,
) raft.ServerAddress {
	if p.serverAddressProvider != nil {
		serverAddressOverride, err := p.serverAddressProvider.ServerAddr(id)
		if err != nil {
			p.logger.Warn("unable to get address for server id %v, using fallback address %v: %v", id, target, err)
		} else {
			return serverAddressOverride
		}
	}
	return target
}

func (p *PgTransport) getConn(target raft.ServerAddress) (*pgConn, error) {
	// Check for a pooled conn
	if conn := p.getPooledConn(target); conn != nil {
		return conn, nil
	}

	// Dial a new connection
	conn, err := p.stream.Dial(target, p.timeout)
	if err != nil {
		p.logger.Error("could not dial connection for [%v]: %v", target, err)
		return nil, err
	}

	wire, err := pgproto.NewRaftWire(conn, conn)
	if err != nil {
		p.logger.Error("could not create wire for [%v]: %v", target, err)
		return nil, err
	}

	pgConn := &pgConn{
		target: target,
		conn:   conn,
		wire:   wire,
	}

	return pgConn, nil
}

func (p *PgTransport) returnConn(conn *pgConn) {
	p.connPoolLock.Lock()
	defer p.connPoolLock.Unlock()

	key := conn.target
	conns, _ := p.connPool[key]

	if !p.IsShutdown() && len(conns) < p.maxPool {
		p.connPool[key] = append(conns, conn)
	} else {
		conn.Release()
	}
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

func (p *PgTransport) listen() {
	for {
		golog.Verbosef("listening for transport connections")
		conn, err := p.stream.Accept()
		if err != nil {
			if p.IsShutdown() {
				return
			}
			p.logger.Error("failed to accept connection: %v", err)
			continue
		}
		p.logger.Debug("%v accepted connection from: %v", p.LocalAddr(), conn.RemoteAddr())

		go p.handleConnection(p.getStreamContext(), conn)
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
			if err != io.EOF {
				p.logger.Error("failed to receive message from [%v]: %v", conn.RemoteAddr(), err)
			}
			return
		}

		responseChannel := make(chan raft.RPCResponse, 1)
		rpc := raft.RPC{
			RespChan: responseChannel,
		}

		isHeartbeat := false
		switch req := request.(type) {
		case *pgproto.AppendEntriesRequest:
			rpc.Command = &req.AppendEntriesRequest

			// Check if this is a heartbeat
			if req.Term != 0 && req.Leader != nil &&
				req.PrevLogEntry == 0 && req.PrevLogTerm == 0 &&
				len(req.Entries) == 0 && req.LeaderCommitIndex == 0 {
				isHeartbeat = true
			}
		case *pgproto.RequestVoteRequest:
			rpc.Command = &req.RequestVoteRequest
		case *pgproto.InstallSnapshotRequest:
			rpc.Command = &req.InstallSnapshotRequest
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
			var msg pgproto.RaftMessage
			switch rsp := response.Response.(type) {
			case *raft.AppendEntriesResponse:
				msg = &pgproto.AppendEntriesResponse{
					Error:                 response.Error,
					AppendEntriesResponse: *rsp,
				}
			case *raft.RequestVoteResponse:
				msg = &pgproto.RequestVoteResponse{
					Error:               response.Error,
					RequestVoteResponse: *rsp,
				}
			case *raft.InstallSnapshotResponse:
				msg = &pgproto.InstallSnapshotResponse{
					Error:                   response.Error,
					InstallSnapshotResponse: *rsp,
				}
			case nil:
				msg = &pgproto.ErrorResponse{
					Message: response.Error.Error(),
				}
			}

			if err := wire.Send(msg); err != nil {
				p.logger.Error("failed to send response to [%v]: %v", conn.RemoteAddr(), err)
				return
			}
		case <-p.shutdownChannel:
			p.logger.Warn("closing transport due to shutdown")
			return
		}
	}
}

func newPgPipeline(trans *PgTransport, conn *pgConn) (*pgPipeline, error) {
	wire, err := pgproto.NewRaftWire(conn.conn, conn.conn)
	if err != nil {
		return nil, err
	}

	p := &pgPipeline{
		conn:              conn,
		transport:         trans,
		doneChannel:       make(chan raft.AppendFuture, rpcMaxPipeline),
		inProgressChannel: make(chan *appendFuture, rpcMaxPipeline),
		shutdownChannel:   make(chan struct{}),
		logger:            logger.NewLogger(),
		wire:              wire,
	}

	go p.processResponses()

	return p, nil
}

// AppendEntries is used to pipeline a new append entries request.
func (p *pgPipeline) AppendEntries(
	args *raft.AppendEntriesRequest,
	resp *raft.AppendEntriesResponse,
) (raft.AppendFuture, error) {
	// Create a new future
	future := &appendFuture{
		start: time.Now(),
		args:  args,
		resp:  resp,
	}
	future.init()

	// Add a send timeout
	if timeout := p.transport.timeout; timeout > 0 {
		p.conn.conn.SetWriteDeadline(time.Now().Add(timeout))
	}

	if err := p.wire.Send(&pgproto.AppendEntriesRequest{
		AppendEntriesRequest: *args,
	}); err != nil {
		return nil, err
	}

	select {
	case p.inProgressChannel <- future:
		return future, nil
	case <-p.shutdownChannel:
		return nil, fmt.Errorf("pipeline shutting down")
	}
}

// Consumer returns a channel that can be used to consume complete futures.
func (p *pgPipeline) Consumer() <-chan raft.AppendFuture {
	return p.doneChannel
}

// Close is used to shutdown the pipeline connection.
func (p *pgPipeline) Close() error {
	p.shutdownLock.Lock()
	defer p.shutdownLock.Unlock()
	if p.shutdown {
		return nil
	}

	// Release the connection
	p.conn.Release()

	p.shutdown = true
	close(p.shutdownChannel)
	return nil
}

func (p *pgPipeline) processResponses() {
	timeout := p.transport.timeout

	for {
		select {
		case future := <-p.inProgressChannel:
			func(future *appendFuture) {
				defer func(future *appendFuture) {
					p.doneChannel <- future
				}(future)

				if timeout > 0 {
					p.conn.conn.SetReadDeadline(time.Now().Add(timeout))
				}

				response, err := p.wire.Receive()
				if err != nil {
					p.logger.Error("could not process response from [%v]: %v", p.conn.conn.RemoteAddr(), err)
					future.respond(err)
					return
				}

				switch msg := response.(type) {
				case *pgproto.AppendEntriesResponse:
					future.resp = &msg.AppendEntriesResponse
					future.respond(msg.Error)
				case *pgproto.ErrorResponse:
					future.resp = nil
					future.respond(errors.New(msg.Message))
				default:
					err = fmt.Errorf("received an invalid response from [%v]: %v", p.conn.conn.RemoteAddr(), msg)
					p.logger.Error(err.Error())
					future.respond(err)
				}
			}(future)
		case <-p.shutdownChannel:
			return
		}
	}
}
