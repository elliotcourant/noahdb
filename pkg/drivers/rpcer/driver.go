package rpcer

import (
	"fmt"
	"github.com/elliotcourant/noahdb/pkg/pgproto"
	"net"
)

type RpcDriver struct {
	id         string
	localAddr  net.Addr
	remoteAddr net.Addr
	conn       net.Conn
	front      pgproto.Frontend
}

func NewRPCDriver(id string, localAddr net.Addr, remoteAddr string) (*RpcDriver, error) {
	driver := &RpcDriver{
		localAddr: localAddr,
		id:        id,
	}
	addr, err := net.ResolveTCPAddr("tcp", remoteAddr)
	if err != nil {
		return nil, err
	}
	driver.remoteAddr = addr
	conn, err := net.Dial("tcp", addr.String())
	if err != nil {
		return nil, err
	}

	driver.conn = conn

	frontend, err := pgproto.NewFrontend(driver.conn, driver.conn)
	if err != nil {
		return nil, err
	}

	if err := frontend.Send(&pgproto.RpcStartupMessage{}); err != nil {
		return nil, err
	}

	response, err := frontend.Receive()
	if err != nil {
		return nil, err
	}

	switch msg := response.(type) {
	case *pgproto.ReadyForQuery:
		driver.front = frontend
		return driver, nil
	case *pgproto.ErrorResponse:
		return nil, fmt.Errorf("could not connect via rpc: %s", msg.Message)
	default:
		return nil, fmt.Errorf("could not handle response message: %v", msg)
	}
}

func (rpc *RpcDriver) Join() error {
	if err := rpc.front.Send(&pgproto.JoinRequest{
		NodeID:  rpc.id,
		Address: rpc.localAddr.String(),
	}); err != nil {
		return err
	}

	response, err := rpc.front.Receive()
	if err != nil {
		return err
	}

	switch msg := response.(type) {
	case *pgproto.ReadyForQuery:
		return nil
	case *pgproto.ErrorResponse:
		return fmt.Errorf("could not join: %s", msg.Message)
	default:
		return fmt.Errorf("could not handle response message when joining: %v", msg)
	}
}

func (rpc *RpcDriver) Discover() (string, error) {
	if err := rpc.front.Send(&pgproto.DiscoveryRequest{}); err != nil {
		return "", err
	}

	response, err := rpc.front.Receive()
	if err != nil {
		return "", err
	}

	switch msg := response.(type) {
	case *pgproto.DiscoveryResponse:
		return msg.LeaderAddr, nil
	case *pgproto.ErrorResponse:
		return "", fmt.Errorf("could not discover: %s", msg.Message)
	default:
		return "", fmt.Errorf("could not handle response message when discovering: %v", msg)
	}
}

func (rpc *RpcDriver) GetSequenceChunk(name string) (*pgproto.SequenceChunkResponse, error) {
	if err := rpc.front.Send(&pgproto.SequenceRequest{
		SequenceName: name,
	}); err != nil {
		return nil, err
	}

	response, err := rpc.front.Receive()
	if err != nil {
		return nil, err
	}

	switch msg := response.(type) {
	case *pgproto.SequenceChunkResponse:
		return msg, nil
	case *pgproto.ErrorResponse:
		return nil, fmt.Errorf("could not discover: %s", msg.Message)
	default:
		return nil, fmt.Errorf("could not handle response message when discovering: %v", msg)
	}
}
