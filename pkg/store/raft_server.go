package store

import (
	"context"
	"fmt"
)

// ClusterServiceServer is the server API for ClusterService service.
type clusterServer struct {
	Store
}

func (server *clusterServer) GetSequenceChunk(ctx context.Context, request *SequenceChunkRequest) (*SequenceChunkResponse, error) {
	return server.getSequenceChunk(request.SequenceName)
}

func (server *clusterServer) SendCommand(ctx context.Context, command *Command) (*CommandResponse, error) {
	switch command.Operation {
	case Operation_DELETE:
		return server.serverDelete(*command)
	case Operation_SET:
		return server.serverSet(*command)
	case Operation_GET:
		return server.serverGet(*command)
	default:
		return nil, fmt.Errorf("could not handle operation %d", command.Operation)
	}
}

func (server *clusterServer) Join(ctx context.Context, join *JoinRequest) (*JoinResponse, error) {
	response := &JoinResponse{}
	if err := server.Store.join(join.ID, join.RaftAddress); err != nil {
		response.IsSuccess = false
		response.ErrorMessage = err.Error()
	} else {
		response.IsSuccess = true
	}
	return response, nil
}

func (server *clusterServer) GetNodeID(ctx context.Context, join *GetNodeIdRequest) (*GetNodeIdResponse, error) {
	nodeId, err := server.NextSequenceValueById("_node_ids_")
	if err != nil {
		return nil, err
	}
	response := &GetNodeIdResponse{
		NodeID: *nodeId,
	}
	return response, nil
}

func (server *clusterServer) serverGet(command Command) (*CommandResponse, error) {
	response := &CommandResponse{
		Operation: command.Operation,
	}
	value, err := server.Store.Get(command.Key)
	if err != nil {
		response.ErrorMessage = err.Error()
		response.IsSuccess = false
	} else {
		response.Value = value
		response.IsSuccess = true
	}
	return response, err
}

func (server *clusterServer) serverSet(command Command) (*CommandResponse, error) {
	response := &CommandResponse{}
	if err := server.Store.Set(command.Key, command.Value); err != nil {
		response.ErrorMessage = err.Error()
		response.IsSuccess = false
	} else {
		response.IsSuccess = true
	}
	response.Operation = command.Operation
	return response, nil
}

func (server *clusterServer) serverDelete(command Command) (*CommandResponse, error) {
	response := &CommandResponse{}
	if err := server.Store.Delete(command.Key); err != nil {
		response.ErrorMessage = err.Error()
		response.IsSuccess = false
	} else {
		response.IsSuccess = true
	}
	response.Operation = command.Operation
	return response, nil
}
