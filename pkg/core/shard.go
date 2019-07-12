package core

import (
	"fmt"
	"github.com/elliotcourant/noahdb/pkg/drivers/rqliter"
	"github.com/elliotcourant/noahdb/pkg/frunk"
	"github.com/elliotcourant/timber"
	"gopkg.in/doug-martin/goqu.v5"
	// Use the postgres adapter for building queries.
	"database/sql"
	_ "gopkg.in/doug-martin/goqu.v5/adapters/postgres"

	_ "github.com/lib/pq"
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

type DataNodePressure struct {
	DataNodeID uint64
	Shards     int
}

type ShardPressure struct {
	ShardID uint64
	Tenants int
}

type shardContext struct {
	*base
}

// ShardContext is just a wrapper interface for shard metadata.
type ShardContext interface {
	NewShard() (Shard, error)
	GetShards() ([]Shard, error)
	GetDataNodeShards() ([]DataNodeShard, error)
	GetWriteDataNodeShards(uint64) ([]DataNodeShard, error)
	BalanceOrphanShards() error
	GetDataNodesPressure(max int) ([]DataNodePressure, error)
	GetShardPressures(max int) ([]ShardPressure, error)
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
	shard.ShardID = id
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
	timber.Debugf("found %d orphaned shards", len(ids))
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
	pressures, err := ctx.GetDataNodesPressure(len(ids))

	for i, shardId := range ids {
		// Determine which node this shard should be assigned to.
		pressureIndex := i % len(pressures)
		dataNode := pressures[pressureIndex]
		timber.Debugf("assigning shard [%d] to data node [%d]", shardId, dataNode.DataNodeID)

		id, err := ctx.db.NextSequenceValueById(dataNodeShardIdSequencePath)
		if err != nil {
			return err
		}

		newDataNodeShard := goqu.
			From("data_node_shards").
			Insert(goqu.Record{
				"data_node_shard_id": id,
				"data_node_id":       dataNode.DataNodeID,
				"shard_id":           shardId,
				"read_only":          false,
			}).Sql
		if _, err := ctx.db.Exec(newDataNodeShard); err != nil {
			return err
		}

		dataNodeMeta, err := ctx.DataNodes().GetDataNode(dataNode.DataNodeID)
		if err != nil {
			timber.Criticalf("failed to retrieve metadata for data node [%d]: %v", dataNode.DataNodeID, err)
			continue
		}

		dataNodeAddress := fmt.Sprintf("%s:%d", dataNodeMeta.GetAddress(), dataNodeMeta.GetPort())
		timber.Debugf("trying to connect to data node [%d] at %s to init shards", dataNode.DataNodeID, dataNodeAddress)

		databaseLogin := dataNodeMeta.GetUser()
		if dataNodeMeta.Password != "" {
			databaseLogin += ":" + dataNodeMeta.GetPassword()
		}

		connStr := fmt.Sprintf("postgres://%s@%s/postgres?sslmode=disable", databaseLogin, dataNodeAddress)
		db, err := sql.Open("postgres", connStr)
		if err != nil {
			timber.Criticalf("failed to connect to data node [%d] address %s: %v", dataNode.DataNodeID, dataNodeAddress, err)
			continue
		}

		dbname := fmt.Sprintf("noahdb_%d", id)

		kickActiveUsers := fmt.Sprintf(`
		SELECT 
			pg_terminate_backend(pg_stat_activity.pid)
		FROM pg_stat_activity
		WHERE pg_stat_activity.datname = '%s'
		AND pid <> pg_backend_pid();`, dbname)
		_, _ = db.Exec(kickActiveUsers)

		deleteExistingShard := fmt.Sprintf("DROP DATABASE IF EXISTS noahdb_%d", id)
		_, err = db.Exec(deleteExistingShard)
		if err != nil {
			timber.Criticalf("could not drop existing shard db: %v", err)
			continue
		}

		createShardQuery := fmt.Sprintf("CREATE DATABASE noahdb_%d", id)
		_, err = db.Exec(createShardQuery)
		if err != nil {
			timber.Criticalf("failed to create data node shard [%d] on data node [%d] address %s: %v", id, dataNode.DataNodeID, dataNodeAddress, err)
			continue
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

func (ctx *shardContext) GetShardPressures(max int) ([]ShardPressure, error) {
	getPressureQuery, _, _ := goqu.
		From("shards").
		Select(
			goqu.I("shards.shard_id"),
			goqu.COUNT(goqu.I("tenants.tenant_id")).As("tenants")).
		LeftJoin(
			goqu.I("tenants"),
			goqu.On(goqu.I("tenants.shard_id").Eq(goqu.I("shards.shard_id")))).
		GroupBy(goqu.I("shards.shard_id")).
		Order(goqu.I("tenants").Asc()).
		Limit(uint(max)).
		ToSql()
	response, err := ctx.db.Query(getPressureQuery)
	if err != nil {
		return nil, err
	}
	rows := rqliter.NewRqlRows(response)
	pressures := make([]ShardPressure, 0)
	for rows.Next() {
		var shardPressure ShardPressure
		if err := rows.Scan(
			&shardPressure.ShardID,
			&shardPressure.Tenants,
		); err != nil {
			return nil, err
		}
		pressures = append(pressures, shardPressure)
	}
	return pressures, nil
}

func (ctx *shardContext) GetDataNodesPressure(max int) ([]DataNodePressure, error) {
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
	response, err := ctx.db.Query(getPressureQuery)
	if err != nil {
		return nil, err
	}
	rows := rqliter.NewRqlRows(response)
	pressures := make([]DataNodePressure, 0)
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
	return pressures, nil
}

func (ctx *shardContext) GetShards() ([]Shard, error) {
	sql, _, _ := getShardsQuery.ToSql()
	response, err := ctx.db.Query(sql)
	if err != nil {
		return nil, err
	}
	return ctx.shardsFromRows(response)
}

func (ctx *shardContext) GetDataNodeShards() ([]DataNodeShard, error) {
	sql, _, _ := getWriteDataNodeShardsQuery.
		ToSql()
	response, err := ctx.db.Query(sql)
	if err != nil {
		return nil, err
	}
	return ctx.dataNodeShardsFromRows(response)
}

func (ctx *shardContext) GetWriteDataNodeShards(id uint64) ([]DataNodeShard, error) {
	sql, _, _ := getWriteDataNodeShardsQuery.
		Where(goqu.Ex{"shard_id": id}).
		Where(goqu.Ex{"read_only": false}).
		ToSql()
	response, err := ctx.db.Query(sql)
	if err != nil {
		return nil, err
	}
	return ctx.dataNodeShardsFromRows(response)
}

func (ctx *shardContext) shardsFromRows(response *frunk.QueryResponse) ([]Shard, error) {
	rows := rqliter.NewRqlRows(response)
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

func (ctx *shardContext) dataNodeShardsFromRows(response *frunk.QueryResponse) ([]DataNodeShard, error) {
	rows := rqliter.NewRqlRows(response)
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
