package core

import (
	"fmt"
	"github.com/elliotcourant/noahdb/pkg/drivers/rqliter"
	"github.com/elliotcourant/noahdb/pkg/frunk"
	"github.com/readystock/goqu"
)

var (
	getDataNodesQuery = goqu.From("data_nodes").
		Select("data_nodes.*")
)

type dataNodeContext struct {
	*base
}

type DataNodeContext interface {
	GetDataNode(id uint64) (DataNode, error)
	GetDataNodes() ([]DataNode, error)
	GetDataNodesForShard(uint64) ([]DataNode, error)
	GetDataNodeForDataNodeShard(uint64) (DataNode, error)
	GetRandomDataNodeShardID() (uint64, error)
	GetDataNodeShardIDs() ([]uint64, error)
	GetDataNodeShardIDsForShard(uint64) ([]uint64, error)
	NewDataNode(address string, port int, user string, password string) (DataNode, error)
}

func (ctx *base) DataNodes() DataNodeContext {
	return &dataNodeContext{
		ctx,
	}
}

func (ctx *dataNodeContext) NewDataNode(address string, port int, user, password string) (DataNode, error) {
	id, err := ctx.db.NextSequenceValueById(dataNodeIdSequencePath)
	if err != nil {
		return DataNode{}, err
	}

	compiledSql := goqu.From("data_nodes").
		Insert(goqu.Record{
			"data_node_id": id,
			"address":      address,
			"port":         port,
			"user":         user,
			"password":     password,
			"healthy":      true,
		}).Sql
	_, err = ctx.db.Exec(compiledSql)
	return DataNode{
		DataNodeID: id,
		Address:    address,
		Port:       int32(port),
		User:       user,
		Password:   password,
		Healthy:    true,
	}, err
}

func (ctx *dataNodeContext) GetDataNode(id uint64) (DataNode, error) {
	compiledQuery, _, _ := getDataNodesQuery.
		Where(goqu.Ex{
			"data_node_id": id,
		}).
		Limit(1).ToSql()
	response, err := ctx.db.Query(compiledQuery)
	if err != nil {
		return DataNode{}, err
	}
	nodes, err := ctx.dataNodesFromRows(response)
	if err != nil {
		return DataNode{}, err
	}
	if len(nodes) == 0 {
		return DataNode{}, fmt.Errorf("could not find data node with ID [%d]", id)
	}
	return nodes[0], nil
}

func (ctx *dataNodeContext) GetDataNodes() ([]DataNode, error) {
	compiledQuery, _, _ := getDataNodesQuery.ToSql()
	response, err := ctx.db.Query(compiledQuery)
	if err != nil {
		return nil, err
	}
	return ctx.dataNodesFromRows(response)
}

func (ctx *dataNodeContext) GetDataNodeShardIDs() ([]uint64, error) {
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

func (ctx *dataNodeContext) GetRandomDataNodeShardID() (uint64, error) {
	compiledQuery, _, _ := goqu.
		From("data_nodes").
		Select("data_node_shards.data_node_shard_id").
		InnerJoin(
			goqu.I("data_node_shards"),
			goqu.On(goqu.I("data_node_shards.data_node_id").Eq(goqu.I("data_nodes.data_node_id")))).
		InnerJoin(
			goqu.I("shards"),
			goqu.On(goqu.I("shards.shard_id").Eq(goqu.I("data_node_shards.shard_id")))).
		Where(goqu.Ex{
			"data_nodes.healthy": true,
			"shards.state":       ShardState_Stable,
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
	if len(nodes) == 0 {
		return DataNode{}, fmt.Errorf("could not find data node for data node shard ID [%d]", id)
	}
	return nodes[0], err
}

func (ctx *dataNodeContext) dataNodesFromRows(response *frunk.QueryResponse) ([]DataNode, error) {
	rows := rqliter.NewRqlRows(response)
	nodes := make([]DataNode, 0)
	for rows.Next() {
		node := DataNode{}
		if err := rows.Scan(
			&node.DataNodeID,
			&node.Address,
			&node.Port,
			&node.User,
			&node.Password,
			&node.Healthy); err != nil {
			return nil, err
		}
		nodes = append(nodes, node)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return nodes, nil
}

func idArray(response *frunk.QueryResponse) ([]uint64, error) {
	rows := rqliter.NewRqlRows(response)
	ids := make([]uint64, 0)
	for rows.Next() {
		id := uint64(0)
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return ids, nil
}

func count(response *frunk.QueryResponse) (uint64, error) {
	rows := rqliter.NewRqlRows(response)
	for rows.Next() {
		count := uint64(0)
		if err := rows.Scan(&count); err != nil {
			return 0, err
		}
		return count, nil
	}
	if err := rows.Err(); err != nil {
		return 0, err
	}
	return 0, nil
}

func exists(response *frunk.QueryResponse) (bool, error) {
	c, err := count(response)
	return c > 0, err
}
