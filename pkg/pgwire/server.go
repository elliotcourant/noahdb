/*
 * Copyright (c) 2019 Ready Stock
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing
 * permissions and limitations under the License.
 */

package pgwire

import (
	"context"
	"fmt"
	"github.com/elliotcourant/noahdb/pkg/pgwire/pgerror"
	"github.com/elliotcourant/noahdb/pkg/pgwire/pgwirebase"
	"github.com/elliotcourant/noahdb/pkg/sql"
	"github.com/elliotcourant/noahdb/pkg/util/syncutil"
	"github.com/pkg/errors"
	"github.com/readystock/golog"
	"github.com/readystock/noah/db/core"
	"io"
	"net"
	"time"
)

const (
	// ErrSSLRequired is returned when a client attempts to connect to a
	// secure server in cleartext.
	ErrSSLRequired = "node is running secure mode, SSL connection required"

	// ErrDraining is returned when a client attempts to connect to a server
	// which is not accepting client connections.
	ErrDraining = "server is not accepting clients"
)

// Fully-qualified names for metrics.

const (
	version30  = 196608
	versionSSL = 80877103
)

// cancelMaxWait is the amount of time a draining server gives to sessions to
// react to cancellation and return before a forceful shutdown.
const cancelMaxWait = 1 * time.Second

// connReservationBatchSize determines for how many connections memory
// is pre-reserved at once.
var connReservationBatchSize = 5

var (
	sslSupported   = []byte{'S'}
	sslUnsupported = []byte{'N'}
)

// cancelChanMap keeps track of channels that are closed after the associated
// cancellation function has been called and the cancellation has taken place.
type cancelChanMap map[chan struct{}]context.CancelFunc

// Server implements the server side of the PostgreSQL wire protocol.
type Server struct {
	// cfg *base.Config
	// execCfg    *sql.ExecutorConfig
	SQLServer *sql.Server

	mu struct {
		syncutil.Mutex
		// connCancelMap entries represent connections started when the server
		// was not draining. Each value is a function that can be called to
		// cancel the associated connection. The corresponding key is a channel
		// that is closed when the connection is done.
		connCancelMap cancelChanMap
		draining      bool
	}
}

// ServerMetrics is the set of metrics for the pgwire server.

// noteworthySQLMemoryUsageBytes is the minimum size tracked by the
// client SQL pool before the pool start explicitly logging overall
// usage growth in the log.

// noteworthyConnMemoryUsageBytes is the minimum size tracked by the
// connection monitor before the monitor start explicitly logging overall
// usage growth in the log.

// MakeServer creates a Server.
//
// Start() needs to be called on the Server so it begins processing.
func MakeServer(
	cfg *base.Config,
) *Server {
	server := &Server{
		cfg: cfg,
	}
	server.SQLServer = sql.NewServer()

	server.mu.Lock()
	server.mu.connCancelMap = make(cancelChanMap)
	server.mu.Unlock()

	return server
}

// Match returns true if rd appears to be a Postgres connection.
func Match(rd io.Reader) bool {
	var buf pgwirebase.ReadBuffer
	_, err := buf.ReadUntypedMsg(rd)
	if err != nil {
		return false
	}
	version, err := buf.GetUint32()
	if err != nil {
		return false
	}
	return version == version30 || version == versionSSL
}

// IsDraining returns true if the server is not currently accepting
// connections.
func (s *Server) IsDraining() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.mu.draining
}

// Drain prevents new connections from being served and waits for drainWait for
// open connections to terminate before canceling them.
// An error will be returned when connections that have been canceled have not
// responded to this cancellation and closed themselves in time. The server
// will remain in draining state, though open connections may continue to
// exist.
// The RFC on drain modes has more information regarding the specifics of
// what will happen to connections in different states:
// https://github.com/cockroachdb/cockroach/blob/master/docs/RFCS/20160425_drain_modes.md
func (s *Server) Drain(drainWait time.Duration) error {
	return s.drainImpl(drainWait, cancelMaxWait)
}

// Undrain switches the server back to the normal mode of operation in which
// connections are accepted.
func (s *Server) Undrain() {
	s.mu.Lock()
	s.setDrainingLocked(false)
	s.mu.Unlock()
}

// setDrainingLocked sets the server's draining state and returns whether the
// state changed (i.e. drain != s.mu.draining). s.mu must be locked.
func (s *Server) setDrainingLocked(drain bool) bool {
	if s.mu.draining == drain {
		return false
	}
	s.mu.draining = drain
	return true
}

func (s *Server) drainImpl(drainWait time.Duration, cancelWait time.Duration) error {
	// This anonymous function returns a copy of s.mu.connCancelMap if there are
	// any active connections to cancel. We will only attempt to cancel
	// connections that were active at the moment the draining switch happened.
	// It is enough to do this because:
	// 1) If no new connections are added to the original map all connections
	// will be canceled.
	// 2) If new connections are added to the original map, it follows that they
	// were added when s.mu.draining = false, thus not requiring cancellation.
	// These connections are not our responsibility and will be handled when the
	// server starts draining again.
	connCancelMap := func() cancelChanMap {
		s.mu.Lock()
		defer s.mu.Unlock()
		if !s.setDrainingLocked(true) {
			// We are already draining.
			return nil
		}
		connCancelMap := make(cancelChanMap)
		for done, cancel := range s.mu.connCancelMap {
			connCancelMap[done] = cancel
		}
		return connCancelMap
	}()
	if len(connCancelMap) == 0 {
		return nil
	}

	// Spin off a goroutine that waits for all connections to signal that they
	// are done and reports it on allConnsDone. The main goroutine signals this
	// goroutine to stop work through quitWaitingForConns.
	allConnsDone := make(chan struct{})
	quitWaitingForConns := make(chan struct{})
	defer close(quitWaitingForConns)
	go func() {
		defer close(allConnsDone)
		for done := range connCancelMap {
			select {
			case <-done:
			case <-quitWaitingForConns:
				return
			}
		}
	}()

	// Wait for all connections to finish up to drainWait.
	select {
	case <-time.After(drainWait):
	case <-allConnsDone:
	}

	// Cancel the contexts of all sessions if the server is still in draining
	// mode.
	if stop := func() bool {
		s.mu.Lock()
		defer s.mu.Unlock()
		if !s.mu.draining {
			return true
		}
		for _, cancel := range connCancelMap {
			// There is a possibility that different calls to SetDraining have
			// overlapping connCancelMaps, but context.CancelFunc calls are
			// idempotent.
			cancel()
		}
		return false
	}(); stop {
		return nil
	}

	select {
	case <-time.After(cancelWait):
		return errors.Errorf("some sessions did not respond to cancellation within %s", cancelWait)
	case <-allConnsDone:
	}
	return nil
}

// ServeConn serves a single connection, driving the handshake process and
// delegating to the appropriate connection type.
func (s *Server) ServeConn(ctx core.Context, conn net.Conn) error { // ctx context.Context,
	// If the Server is draining, we will use the connection only to send an
	// error, so we don't count it in the stats. This makes sense since
	// DrainClient() waits for that number to drop to zero,
	// so we don't want it to oscillate unnecessarily.

	var buf pgwirebase.ReadBuffer
	_, err := buf.ReadUntypedMsg(conn)
	if err != nil {
		return err
	}
	version, err := buf.GetUint32()
	if err != nil {
		return err
	}
	errSSLRequired := false
	// if version == versionSSL {
	//     if len(buf.Msg) > 0 {
	//         return errors.Errorf("unexpected data after SSLRequest: %q", buf.Msg)
	//     }
	//
	//     if s.cfg.Insecure {
	//         if _, err := conn.Write(sslUnsupported); err != nil {
	//             return err
	//         }
	//     } else {
	//         if _, err := conn.Write(sslSupported); err != nil {
	//             return err
	//         }
	//         tlsConfig, err := s.cfg.GetServerTLSConfig()
	//         if err != nil {
	//             return err
	//         }
	//         conn = tls.Server(conn, tlsConfig)
	//     }
	//
	//     _, err := buf.ReadUntypedMsg(conn)
	//     if err != nil {
	//         return err
	//     }
	//     version, err = buf.GetUint32()
	//     if err != nil {
	//         return err
	//     }
	// } else if !s.cfg.Insecure {
	//     errSSLRequired = true
	// }

	sendErr := func(err error) error {
		msgBuilder := newWriteBuffer()
		_ /* err */ = writeErr(err, msgBuilder, conn)
		_ = conn.Close()
		return err
	}

	if version != version30 {
		return sendErr(fmt.Errorf("unknown protocol version %d", version))
	}
	if errSSLRequired {
		return sendErr(pgerror.NewError(pgerror.CodeProtocolViolationError, ErrSSLRequired))
	}
	// if draining {
	//     return sendErr(newAdminShutdownErr(errors.New(ErrDraining)))
	// }

	var sArgs sql.SessionArgs
	if sArgs, err = parseOptions(buf.Msg); err != nil {
		return sendErr(pgerror.NewError(pgerror.CodeProtocolViolationError, err.Error()))
	}
	// sArgs.User = tree.Name(sArgs.User).Normalize()

	// Reserve some memory for this connection using the server's monitor. This
	// reduces pressure on the shared pool because the server monitor allocates in
	// chunks from the shared pool and these chunks should be larger than
	// baseSQLMemoryBudget.
	return serveConn(conn, sArgs, s.SQLServer, s.cfg.Insecure, ctx)
}

func parseOptions(data []byte) (sql.SessionArgs, error) {
	args := sql.SessionArgs{}
	buf := pgwirebase.ReadBuffer{Msg: data}
	for {
		key, err := buf.GetString()
		if err != nil {
			return sql.SessionArgs{}, errors.Errorf("error reading option key: %s", err)
		}
		if len(key) == 0 {
			break
		}
		value, err := buf.GetString()
		if err != nil {
			return sql.SessionArgs{}, errors.Errorf("error reading option value: %s", err)
		}
		switch key {
		case "database":
			args.Database = value
		case "user":
			args.User = value
		case "application_name":
			args.ApplicationName = value
		default:
			golog.Warnf("Unrecognized configuration parameter %s", key)
		}
	}
	return args, nil
}

func newAdminShutdownErr(err error) error {
	return pgerror.NewErrorf(pgerror.CodeAdminShutdownError, err.Error())
}
