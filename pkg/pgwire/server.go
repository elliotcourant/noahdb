package pgwire

import (
	"fmt"
	"github.com/elliotcourant/noahdb/pkg/commands"
	"github.com/elliotcourant/noahdb/pkg/core"
	"github.com/elliotcourant/noahdb/pkg/pgproto"
	"github.com/elliotcourant/noahdb/pkg/sql"
	"github.com/elliotcourant/noahdb/pkg/util/stmtbuf"
	"github.com/readystock/golog"
	"net"
	"reflect"
	"strings"
)

type TransportWrapper interface {
	NormalTransport() net.Listener
	ForwardToRaft(net.Conn, error)
	ForwardToRpc(net.Conn, error)
	Close()
}

type ServerConfig interface {
	Address() string
	Port() int
}

func NewServer(colony core.Colony, transport TransportWrapper) error {
	defer transport.Close()

	ln := transport.NormalTransport()

	for {
		golog.Infof("accepting connection at: %s", ln.Addr())
		conn, err := ln.Accept()
		if err != nil {
			return err
		}
		golog.Infof("accepted connection from: %s", conn.RemoteAddr())

		go func() {
			wire, err := newWire(colony, conn)
			if err != nil {
				golog.Errorf("failed setting up wire: %s", err.Error())
			}

			if err := wire.Serve(transport); err != nil {
				golog.Errorf("failed serving connection: %s", err.Error())
			}
		}()
	}
}

type wireServer struct {
	colony core.Colony
	conn   net.Conn

	backend *pgproto.Backend

	stmtBuf stmtbuf.StatementBuffer
}

func newWire(colony core.Colony, conn net.Conn) (*wireServer, error) {
	backend, err := pgproto.NewBackend(conn, conn)
	if err != nil {
		return nil, err
	}
	return &wireServer{
		colony:  colony,
		conn:    conn,
		backend: backend,
	}, nil
}

func (wire *wireServer) Serve(wrapper TransportWrapper) error {
	// Receive startup messages.
	startupMsg, err := wire.backend.ReceiveStartupMessage()
	if err != nil {
		switch err {
		case pgproto.RaftStartupMessageError:
			wrapper.ForwardToRaft(wire.conn, nil)
			return nil
		case pgproto.RpcStartupMessageError:
			wrapper.ForwardToRpc(wire.conn, nil)
			return nil
		default:
			defer wire.conn.Close()
			return wire.Errorf(err.Error())
		}
	}
	defer wire.conn.Close()

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
		if err := sql.Run(wire, terminateChannel); err != nil {
			golog.Errorf(err.Error())
		}
	}()

	for {
		message, err := wire.backend.Receive()
		if err != nil {
			return wire.Errorf(err.Error())
		}

		switch msg := message.(type) {
		case *pgproto.Query:
			if err := wire.handleSimpleQuery(msg); err != nil {
				return wire.StatementBuffer().Push(commands.SendError{
					Err: err,
				})
			}
			if err := wire.StatementBuffer().Push(commands.Sync{}); err != nil {
				return wire.StatementBuffer().Push(commands.SendError{
					Err: err,
				})
			}
		case *pgproto.Execute:
		case *pgproto.Parse:
			if err := wire.handleParse(msg); err != nil {
				return wire.StatementBuffer().Push(commands.SendError{
					Err: err,
				})
			}
		case *pgproto.Describe:
		case *pgproto.Bind:
		case *pgproto.Close:
		case *pgproto.Terminate:
			terminateChannel <- true
			return nil
		case *pgproto.Sync:
		case *pgproto.Flush:
		case *pgproto.CopyData:
		default:
			return wire.Errorf("could not handle message type [%s]", reflect.TypeOf(message).Elem().Name())
		}
	}
	return nil
}

func (wire *wireServer) Backend() *pgproto.Backend {
	return wire.backend
}

func (wire *wireServer) StatementBuffer() stmtbuf.StatementBuffer {
	return wire.stmtBuf
}

func (wire *wireServer) Colony() core.Colony {
	return wire.colony
}

func (wire *wireServer) Errorf(message string, args ...interface{}) error {
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
