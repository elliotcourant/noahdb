package rpcwire

import (
	"fmt"
	"github.com/elliotcourant/noahdb/pkg/pgproto"
	"github.com/readystock/golog"
)

func (wire *rpcWire) handleJoin(join *pgproto.JoinRequest) error {
	golog.Infof("received join request from node ID [%s]", join.NodeID)
	if !wire.colony.IsLeader() {
		golog.Warnf("join request needs to be forwarded to leader")
		return fmt.Errorf("cannot handle join request, node is not leader")
	}
	if err := wire.colony.Join(join.NodeID, join.Address); err != nil {
		return err
	} else {
		return wire.backend.Send(&pgproto.ReadyForQuery{})
	}
}
