package core

import (
	"database/sql"
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
	GetRandomDataNode() (DataNode, error)
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
	rows, err := ctx.db.Query(compiledQuery)
	if err != nil {
		return nil, err
	}
	return ctx.dataNodesFromRows(rows)
}

func (ctx *dataNodeContext) GetRandomDataNode() (DataNode, error) {
	compiledQuery, _, _ := getDataNodesQuery.
		Where(goqu.Ex{
			"healthy": true,
		}).
		Order(goqu.L("RANDOM()").Asc()).
		Limit(1).
		ToSql()
	rows, err := ctx.db.Query(compiledQuery)
	if err != nil {
		return DataNode{}, err
	}
	nodes, err := ctx.dataNodesFromRows(rows)
	return nodes[0], err
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
	rows, err := ctx.db.Query(compiledQuery)
	if err != nil {
		return nil, err
	}
	return ctx.dataNodesFromRows(rows)
}

func (ctx *dataNodeContext) dataNodesFromRows(rows *sql.Rows) ([]DataNode, error) {
	defer rows.Close()
	nodes := make([]DataNode, 0)
	for rows.Next() {
		node := DataNode{}
		if err := rows.Scan(
			&node.DataNodeID,
			&node.Address,
			&node.Port,
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
