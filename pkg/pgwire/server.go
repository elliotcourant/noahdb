package pgwire

import (
	"fmt"
	"github.com/elliotcourant/noahdb/pkg/pgwire/pgproto"
	"github.com/readystock/golog"
	"io"
	"net"
	"reflect"
	"strings"
)

type ServerConfig interface {
	Address() string
	Port() int
}

func NewServer(config ServerConfig) error {
	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", config.Address(), config.Port()))
	if err != nil {
		return err
	}
	defer ln.Close()

	golog.Verbosef("listening for connections on: %s:%d", config.Address(), config.Port())

	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}

		go func() {
			defer conn.Close()
			wire, err := newWire(conn)
			if err != nil {
				golog.Errorf("failed setting up wire: %s", err.Error())
			}

			golog.Verbosef("handling connection from: %s", conn.RemoteAddr().String())
			if err := wire.Serve(); err != nil {
				golog.Errorf("failed serving connection: %s", err.Error())
			}
		}()
	}
}

type wireServer struct {
	reader  io.Reader
	writer  io.Writer
	backend *pgproto.Backend
}

func newWire(conn net.Conn) (*wireServer, error) {
	backend, err := pgproto.NewBackend(conn, conn)
	if err != nil {
		return nil, err
	}
	return &wireServer{
		reader:  conn,
		writer:  conn,
		backend: backend,
	}, nil
}

func (wire *wireServer) Serve() error {
	// Receive startup messages.
	startupMsg, err := wire.backend.ReceiveStartupMessage()
	if err != nil {
		return wire.Errorf(err.Error())
	}
	golog.Infof("received startup message: %+v", startupMsg)

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

		authResponse, ok := response.(*pgproto.PasswordMessage)
		if !ok {
			return wire.Errorf("authentication failed")
		}
		golog.Verbosef("received password authentication for user [%s]", authResponse.Password)

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

	for {
		message, err := wire.backend.Receive()
		if err != nil {
			return wire.Errorf(err.Error())
		}

		switch message.(type) {
		case *pgproto.Query:
			if err := wire.backend.Send(&pgproto.CommandComplete{
				CommandTag: "SELECT 1",
			}); err != nil {
				return wire.Errorf(err.Error())
			}
		case *pgproto.Execute:
		case *pgproto.Parse:
		case *pgproto.Describe:
		case *pgproto.Bind:
		case *pgproto.Close:
		case *pgproto.Terminate:
		case *pgproto.Sync:
		case *pgproto.Flush:
		case *pgproto.CopyData:
		default:
			return wire.Errorf("could not handle message type [%s]", reflect.TypeOf(message).Elem().Name())
		}

		if err := wire.backend.Send(&pgproto.ReadyForQuery{
			TxStatus: 'I',
		}); err != nil {
			return wire.Errorf(err.Error())
		}
	}
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
