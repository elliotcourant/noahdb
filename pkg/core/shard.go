package core

import (
	"database/sql"
	"fmt"
	"gopkg.in/doug-martin/goqu.v5"
	_ "gopkg.in/doug-martin/goqu.v5/adapters/postgres"
)

var (
	getShardsQuery = goqu.
			From("shards").
			Select(
			"shard_id",
			"ready")
	getWriteDataNodeShardsQuery = goqu.
					From("data_node_shards").
					Select(
			"data_node_shard_id",
			"data_node_id",
			"shard_id",
			"read_only")
)

type shardContext struct {
	*base
}

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
	sql := fmt.Sprintf("INSERT INTO shards VALUES (%d, %d);", shard.ShardID, shard.State)
	_, err = ctx.db.Exec(sql)
	if err != nil {
		return Shard{}, err
	}
	return shard, nil
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
