package core

import (
	"database/sql"
	"github.com/readystock/golog"
	"gopkg.in/doug-martin/goqu.v5"
	// Use the postgres adapter for building queries.
	_ "gopkg.in/doug-martin/goqu.v5/adapters/postgres"
)

var (
	getShardsQuery = goqu.
			From("shards").
			Select(
			"shard_id",
			"state")
	getWriteDataNodeShardsQuery = goqu.
					From("data_node_shards").
					Select(
			"data_node_shard_id",
			"data_node_id",
			"shard_id",
			"read_only")
)

type dataNodePressure struct {
	DataNodeID uint64
	Shards     int
}

type shardContext struct {
	*base
}

// ShardContext is just a wrapper interface for shard metadata.
type ShardContext interface {
	NewShard() (Shard, error)
	GetShards() ([]Shard, error)
	GetWriteDataNodeShards(uint64) ([]DataNodeShard, error)
}

func (ctx *base) Shards() ShardContext {
	return &shardContext{
		ctx,
	}
}

func (ctx *shardContext) NewShard() (shard Shard, err error) {
	id, err := ctx.db.NextSequenceValueById(shardIdSequencePath)
	if err != nil {
		return shard, err
	}
	shard.ShardID = *id
	shard.State = ShardState_New

	compiledQuery := goqu.
		From("shards").
		Insert(goqu.Record{
			"shard_id": shard.ShardID,
			"state":    shard.State,
		}).Sql

	_, err = ctx.db.Exec(compiledQuery)
	if err != nil {
		return Shard{}, err
	}
	return shard, nil
}

// BalanceOrphanShards looks at all of the shards in the cluster
// that are not currently associated with a data node and assigns
// them a data node. Then marks that shard as ready.
func (ctx *shardContext) BalanceOrphanShards() error {
	orphanedShardsQuery, _, _ := goqu.
		From("shards").
		Select("shards.shard_id").
		LeftJoin(
			goqu.I("data_node_shards"),
			goqu.On(goqu.I("data_node_shards.shard_id").Eq(goqu.I("shards.shard_id")))).
		Where(goqu.Ex{
			"data_node_shards.shard_id": nil,
		}).
		ToSql()
	rows, err := ctx.db.Query(orphanedShardsQuery)
	if err != nil {
		return err
	}
	ids, err := idArray(rows)
	if err != nil {
		return err
	}
	golog.Debugf("found %d orphaned shards", len(ids))
	updateShardStateQuery := goqu.
		From("shards").
		Where(goqu.Ex{
			"shard_id": ids,
		}).
		Update(goqu.Ex{
			"state": ShardState_Balancing,
		}).Sql
	_, err = ctx.db.Exec(updateShardStateQuery)
	if err != nil {
		return err
	}
	pressures, err := ctx.getDataNodesPressure(len(ids))

	for i, shardId := range ids {
		// Determine which node this shard should be assigned to.
		pressureIndex := i % len(pressures)
		dataNode := pressures[pressureIndex]
		golog.Debugf("assigning shard [%d] to data node [%d]", shardId, dataNode.DataNodeID)

		id, err := ctx.db.NextSequenceValueById(dataNodeShardIdSequencePath)
		if err != nil {
			return err
		}

		newDataNodeShard := goqu.
			From("data_node_shards").
			Insert(goqu.Record{
				"data_node_shard_id": *id,
				"data_node_id":       dataNode.DataNodeID,
				"shard_id":           shardId,
				"read_only":          false,
			}).Sql
		if _, err := ctx.db.Exec(newDataNodeShard); err != nil {
			return err
		}

		updateShardStatus := goqu.
			From("shards").
			Where(goqu.Ex{
				"shard_id": shardId,
			}).
			Update(goqu.Ex{
				"state": ShardState_Stable,
			}).Sql
		if _, err := ctx.db.Exec(updateShardStatus); err != nil {
			return err
		}
	}
	return nil
}

func (ctx *shardContext) getDataNodesPressure(max int) ([]dataNodePressure, error) {
	getPressureQuery, _, _ := goqu.
		From("data_nodes").
		Select(
			goqu.I("data_nodes.data_node_id"),
			goqu.COUNT(goqu.I("data_node_shards.shard_id")).As("shards")).
		LeftJoin(
			goqu.I("data_node_shards"),
			goqu.On(goqu.I("data_node_shards.data_node_id").Eq(goqu.I("data_nodes.data_node_id")))).
		GroupBy(goqu.I("data_nodes.data_node_id")).
		Order(goqu.I("shards").Asc()).
		Limit(uint(max)).
		ToSql()
	rows, err := ctx.db.Query(getPressureQuery)
	if err != nil {
		return nil, err
	}
	pressures := make([]dataNodePressure, 0)
	for rows.Next() {
		dataNodeId, shards := uint64(0), 0
		if err := rows.Scan(&dataNodeId, &shards); err != nil {
			return nil, err
		}
		pressures = append(pressures, struct {
			DataNodeID uint64
			Shards     int
		}{DataNodeID: dataNodeId, Shards: shards})
	}
	rows.Close()
	return pressures, nil
}

func (ctx *shardContext) GetShards() ([]Shard, error) {
	sql, _, _ := getShardsQuery.ToSql()
	rows, err := ctx.db.Query(sql)
	if err != nil {
		return nil, err
	}
	return ctx.shardsFromRows(rows)
}

func (ctx *shardContext) GetWriteDataNodeShards(id uint64) ([]DataNodeShard, error) {
	sql, _, _ := getWriteDataNodeShardsQuery.
		Where(goqu.Ex{"shard_id": id}).
		Where(goqu.Ex{"read_only": false}).
		ToSql()
	rows, err := ctx.db.Query(sql)
	if err != nil {
		return nil, err
	}
	return ctx.dataNodeShardsFromRows(rows)
}

func (ctx *shardContext) shardsFromRows(rows *sql.Rows) ([]Shard, error) {
	defer rows.Close()
	shards := make([]Shard, 0)
	for rows.Next() {
		shard := Shard{}
		if err := rows.Scan(
			&shard.ShardID, &shard.State); err != nil {
			return nil, err
		}
		shards = append(shards, shard)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return shards, nil
}

func (ctx *shardContext) dataNodeShardsFromRows(rows *sql.Rows) ([]DataNodeShard, error) {
	defer rows.Close()
	dataNodeShards := make([]DataNodeShard, 0)
	for rows.Next() {
		dataNodeShard := DataNodeShard{}
		if err := rows.Scan(
			&dataNodeShard.DataNodeShardID,
			&dataNodeShard.DataNodeID,
			&dataNodeShard.ShardID,
			&dataNodeShard.ReadOnly); err != nil {
			return nil, err
		}
		dataNodeShards = append(dataNodeShards, dataNodeShard)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return dataNodeShards, nil
}
