package store

import (
	"context"
	"fmt"
	"github.com/kataras/go-errors"
	"github.com/kataras/golog"
	"github.com/readystock/raft"
	"google.golang.org/grpc"
	"sync"
)

type clusterClient struct {
	Store
	conn    *grpc.ClientConn
	sync    *sync.Mutex
	addr    raft.ServerAddress
	cluster *clusterServiceClient
}

func (client *clusterClient) validateConnection(leaderAddr raft.ServerAddress) error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("recovered from ", r)
		}
	}()
	client.sync.Lock()
	defer client.sync.Unlock()
	if !client.Store.IsLeader() {
		// If the address is not the same (the leader has changed) then update the connection and reconnect.
		if client.conn != nil {
			client.conn.Close()
		}
		golog.Debugf("[%d] Connecting to leader at `%s`", client.Store.nodeId, leaderAddr)
		conn, err := grpc.Dial(string(leaderAddr), grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			return err
		}
		golog.Debugf("[%d] Connected to leader", client.Store.nodeId)
		client.cluster = &clusterServiceClient{cc: conn}
	}
	return nil
}

func (client *clusterClient) sendCommand(command *Command) (*CommandResponse, error) {
	if err := client.validateConnection(client.Store.raft.Leader()); err != nil {
		return nil, err
	}
	golog.Debugf("[%d] Sending command to leader", client.Store.nodeId)
	if result, err := client.cluster.SendCommand(context.Background(), command); err != nil {
		return nil, err
	} else if !result.IsSuccess {
		return nil, errors.New(result.ErrorMessage)
	} else {
		return result, nil
	}
}

func (client *clusterClient) getNextChunkInSequence(sequenceName string) (*SequenceChunkResponse, error) {
	if err := client.validateConnection(client.Store.raft.Leader()); err != nil {
		return nil, err
	}
	if result, err := client.cluster.GetSequenceChunk(context.Background(), &SequenceChunkRequest{SequenceName: sequenceName}); err != nil {
		return nil, err
	} else {
		return result, nil
	}
}

func (client *clusterClient) joinCluster(addr string) (*JoinResponse, error) {
	return client.cluster.Join(context.Background(), &JoinRequest{RaftAddress: client.listen, ID: client.nodeId})
}
