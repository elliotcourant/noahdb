package core

import (
	"database/sql"
	"github.com/elliotcourant/noahdb/pkg/frunk"
	"github.com/readystock/goqu"
	"strconv"
)

var (
	getDataNodesQuery = goqu.From("data_nodes").
		Select("data_nodes.*")
)

type dataNodeContext struct {
	*base
}

type DataNodeContext interface {
	GetDataNodes() ([]DataNode, error)
	GetDataNodesForShard(uint64) ([]DataNode, error)
	GetDataNodeForDataNodeShard(uint64) (DataNode, error)
	GetRandomDataNodeShardID() (uint64, error)
	GetDataNodeShardIDsForShard(uint64) ([]uint64, error)
	NewDataNode(string, string, string) (DataNode, error)
}

func (ctx *base) DataNodes() DataNodeContext {
	return &dataNodeContext{
		ctx,
	}
}

func (ctx *dataNodeContext) NewDataNode(address, password, port string) (DataNode, error) {
	id, err := ctx.db.NextSequenceValueById(dataNodeIdSequencePath)
	if err != nil {
		return DataNode{}, err
	}

	portInt, err := strconv.Atoi(port)
	if err != nil {
		return DataNode{}, err
	}

	compiledSql := goqu.From("data_nodes").
		Insert(goqu.Record{
			"data_node_id": *id,
			"address":      address,
			"port":         portInt,
			"healthy":      true,
		}).Sql
	_, err = ctx.db.Exec(compiledSql)
	return DataNode{
		DataNodeID: *id,
		Address:    address,
		Port:       int32(portInt),
		Healthy:    true,
	}, err
}

func (ctx *dataNodeContext) GetDataNodes() ([]DataNode, error) {
	compiledQuery, _, _ := getDataNodesQuery.ToSql()
	result, err := ctx.db.Query(compiledQuery)
	if err != nil {
		return nil, err
	}
	return ctx.dataNodesFromRows(result)
}

func (ctx *dataNodeContext) GetRandomDataNodeShardID() (uint64, error) {
	compiledQuery, _, _ := goqu.
		From("data_nodes").
		Select("data_node_shards.data_node_shard_id").
		InnerJoin(
			goqu.I("data_node_shards"),
			goqu.On(goqu.I("data_node_shards.data_node_id").Eq(goqu.I("data_nodes.data_node_id")))).
		Where(goqu.Ex{
			"data_nodes.healthy": true,
		}).
		Order(goqu.L("RANDOM()").Asc()).
		Limit(1).
		ToSql()
	response, err := ctx.db.Query(compiledQuery)
	if err != nil {
		return 0, err
	}
	return count(response)
}

func (ctx *dataNodeContext) GetDataNodesForShard(id uint64) ([]DataNode, error) {
	compiledQuery, _, _ := getDataNodesQuery.
		InnerJoin(
			goqu.I("data_node_shards"),
			goqu.On(goqu.I("data_node_shards.data_node_id").Eq(goqu.I("data_nodes.data_node_id")))).
		Where(goqu.Ex{
			"data_node_shards.shard_id": id,
		}).
		ToSql()
	response, err := ctx.db.Query(compiledQuery)
	if err != nil {
		return nil, err
	}
	return ctx.dataNodesFromRows(response)
}

func (ctx *dataNodeContext) GetDataNodeShardIDsForShard(id uint64) ([]uint64, error) {
	compiledQuery, _, _ := goqu.
		From("data_nodes").
		Select("data_node_shards.data_node_shard_id").
		InnerJoin(
			goqu.I("data_node_shards"),
			goqu.On(goqu.I("data_node_shards.data_node_id").Eq(goqu.I("data_nodes.data_node_id")))).
		Where(goqu.Ex{
			"data_nodes.healthy": true,
		}).
		ToSql()
	response, err := ctx.db.Query(compiledQuery)
	if err != nil {
		return nil, err
	}
	return idArray(response)
}

func (ctx *dataNodeContext) GetDataNodeForDataNodeShard(id uint64) (DataNode, error) {
	compiledQuery, _, _ := getDataNodesQuery.
		InnerJoin(
			goqu.I("data_node_shards"),
			goqu.On(goqu.I("data_node_shards.data_node_id").Eq(goqu.I("data_nodes.data_node_id")))).
		Where(goqu.Ex{
			"data_node_shards.data_node_shard_id": id,
		}).
		ToSql()
	response, err := ctx.db.Query(compiledQuery)
	if err != nil {
		return DataNode{}, err
	}
	nodes, err := ctx.dataNodesFromRows(response)
	return nodes[0], err
}

func (ctx *dataNodeContext) dataNodesFromRows(response *frunk.QueryResponse) ([]DataNode, error) {
	nodes := make([]DataNode, 0)
	for _, row := range response.Rows[0].Values {
		nodes = append(nodes, DataNode{
			DataNodeID: row[0].(uint64),
			Address:    row[1].(string),
			Port:       row[2].(int32),
			Healthy:    row[3].(bool),
		})
	}
	return nodes, nil
}

func count(response *frunk.QueryResponse) (uint64, error) {
	return response.Rows[0].Values[0][0].(uint64), nil
}

func idArray(response *frunk.QueryResponse) ([]uint64, error) {
	ids := make([]uint64, 0)
	for _, row := range response.Rows[0].Values {
		ids = append(ids, row[0].(uint64))
	}
	return ids, nil
}
