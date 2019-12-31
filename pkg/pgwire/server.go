package pgwire

import (
	"fmt"
	"github.com/elliotcourant/noahdb/pkg/commands"
	"github.com/elliotcourant/noahdb/pkg/engine"
	"github.com/elliotcourant/noahdb/pkg/pgproto"
	"github.com/elliotcourant/noahdb/pkg/sql"
	"github.com/elliotcourant/noahdb/pkg/util/stmtbuf"
	"github.com/elliotcourant/timber"
	"github.com/readystock/golog"
	"io"
	"net"
	"reflect"
	"strings"
)

var (
	_ sql.Server = &Server{}
)

type (
	TransportWrapper interface {
		NormalTransport() net.Listener
		ForwardToRaft(net.Conn, error)
		ForwardToRpc(net.Conn, error)
		Close()
	}

	Server struct {
		core        engine.Core
		transaction engine.Transaction
		backend     *pgproto.Backend
		stmtBuf     stmtbuf.StatementBuffer
		log         timber.Logger
	}
)

func RunPgServer(colony engine.Core, listener net.Listener) error {
	defer listener.Close()

	for {
		timber.Verbosef("accepting connection at: %s", listener.Addr())
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		timber.Verbosef("accepted connection from: %s", conn.RemoteAddr())

		go func(conn net.Conn) {
			log := timber.New().Prefix(conn.RemoteAddr().String())
			defer func() {
				if err := conn.Close(); err != nil {
					log.Warningf("error returned when closing connection: %v", err)
				}
			}()
			wire, err := newServer(colony, conn, conn, log)
			if err != nil {
				log.Errorf("failed setting up wire: %s", err.Error())
			}

			if wire == nil {
				log.Errorf("wire is null, cannot continue")
				return
			}

			// Receive startup messages.
			startupMsg, err := wire.backend.ReceiveStartupMessage()
			if err != nil {
				log.Errorf("error receiving startup message: %v", err)
				return
			}
			if err := wire.Serve(*startupMsg); err != nil {
				log.Errorf("failed serving connection: %s", err.Error())
			}
		}(conn)
	}
}

func newServer(core engine.Core, reader io.Reader, writer io.Writer, logger timber.Logger) (*Server, error) {
	backend, err := pgproto.NewBackend(reader, writer)
	if err != nil {
		return nil, err
	}
	return &Server{
		core:    core,
		backend: backend,
		log:     logger,
	}, nil
}

func (wire *Server) Serve(startupMsg pgproto.StartupMessage) error {
	wire.stmtBuf = stmtbuf.NewStatementBuffer() // We only want to setup a statement buffer if there is a need
	if user, ok := startupMsg.Parameters["user"]; !ok || strings.TrimSpace(user) == "" {
		return wire.Errorf("user authentication required")
	} else if username := strings.ToLower(strings.TrimSpace(user)); username != "noah" {
		return wire.Errorf("user [%s] does not exist", username)
	}

	switch startupMsg.ProtocolVersion {
	case pgproto.ProtocolVersionNumber:
		if err := wire.backend.Send(&pgproto.Authentication{
			Type: pgproto.AuthTypeMD5Password,
		}); err != nil {
			return wire.Errorf(err.Error())
		}

		response, err := wire.backend.Receive()
		if err != nil {
			return wire.Errorf(err.Error())
		}

		_, ok := response.(*pgproto.PasswordMessage)
		if !ok {
			return wire.Errorf("authentication failed")
		}

		if err := wire.backend.Send(&pgproto.Authentication{
			Type: pgproto.AuthTypeOk,
		}); err != nil {
			return wire.Errorf(err.Error())
		}

		if err := wire.backend.Send(&pgproto.BackendKeyData{
			ProcessID: 0,
			SecretKey: 0,
		}); err != nil {
			return wire.Errorf(err.Error())
		}

		if err := wire.backend.Send(&pgproto.ReadyForQuery{
			TxStatus: 'I',
		}); err != nil {
			return wire.Errorf(err.Error())
		}
	default:
		return wire.Errorf("could not handle protocol version [%d]", startupMsg.ProtocolVersion)
	}

	terminateChannel := make(chan bool)

	go func() {
		if err := sql.Run(wire, wire.log, terminateChannel); err != nil {
			golog.Errorf(err.Error())
		}
	}()

	for {
		message, err := wire.backend.Receive()
		if err != nil {
			return wire.Errorf(err.Error())
		}

		err = func(message pgproto.FrontendMessage) error {
			switch msg := message.(type) {
			case *pgproto.Query:
				if err := wire.handleSimpleQuery(msg); err != nil {
					return err
				}
				return wire.StatementBuffer().Push(commands.Sync{})
			case *pgproto.Execute:
				return wire.StatementBuffer().Push(commands.ExecutePortal{
					Name:  msg.Portal,
					Limit: int(msg.MaxRows),
				})
			case *pgproto.Parse:
				return wire.handleParse(msg)
			case *pgproto.Describe:
				return wire.handleDescribe(msg)
			case *pgproto.Bind:
				return wire.handleBind(msg)
			case *pgproto.Close:
			case *pgproto.Terminate:
				terminateChannel <- true
				return nil
			case *pgproto.Sync:
				return wire.stmtBuf.Push(commands.Sync{})
			case *pgproto.Flush:
			case *pgproto.CopyData:
			default:
				return wire.Errorf("could not handle message type [%s]", reflect.TypeOf(message).Elem().Name())
			}

			return nil
		}(message)

		if err != nil {
			if e := wire.StatementBuffer().Push(commands.SendError{
				Err: err,
			}); e != nil {
				return e
			}

			if e := wire.StatementBuffer().Push(commands.Sync{}); e != nil {
				return e
			}
		}
	}
}

func (wire *Server) Commit() error {
	if err := wire.transaction.Commit(); err != nil {
		return err
	}

	txn, err := wire.core.Begin()
	wire.transaction = txn

	return err
}

func (wire *Server) Rollback() error {
	if err := wire.transaction.Rollback(); err != nil {
		return err
	}

	txn, err := wire.core.Begin()
	wire.transaction = txn

	return err
}

func (wire *Server) Transaction() engine.Transaction {
	return wire.transaction
}

func (wire *Server) Backend() *pgproto.Backend {
	return wire.backend
}

func (wire *Server) StatementBuffer() stmtbuf.StatementBuffer {
	return wire.stmtBuf
}

func (wire *Server) Errorf(message string, args ...interface{}) error {
	errorMessage := &pgproto.ErrorResponse{
		Severity: "FATAL",
		Code:     "0000",
		Message:  fmt.Sprintf(message, args...),
	}
	if err := wire.backend.Send(errorMessage); err != nil {
		return err
	}
	return fmt.Errorf(message, args...)
}
