package executor

import (
	"github.com/elliotcourant/noahdb/pkg/pgproto"
)

// Tunnel is an interface used to receive and handle
// result data from the data node shard.
// When a query is executed, the response messages
// are piped into the provided tunnel.
type Tunnel interface {
	SendMessage(msg pgproto.BackendMessage) error
}

// ClientTunnel is an implementation of Tunnel
// that will return results from the cluster to
// the end client.
type ClientTunnel struct {
	conn               *pgproto.Backend
	sentRowDescription bool
}

// NewClientTunnel creates a new client to return results to a
// client.
func NewClientTunnel(clientConn *pgproto.Backend, extendedQueryMode bool) *ClientTunnel {
	return &ClientTunnel{
		conn:               clientConn,
		sentRowDescription: extendedQueryMode,
	}
}

// SendMessage will take the provided message and pass it to the client.
// It will only send row descriptions once, and only if the tunnel
// is not in extended query mode.
func (t *ClientTunnel) SendMessage(msg pgproto.BackendMessage) error {
	switch msg.(type) {
	case *pgproto.RowDescription:
		if t.sentRowDescription {
			return nil
		}
		t.sentRowDescription = true
	}
	return t.conn.Send(msg)
}
