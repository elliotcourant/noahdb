package core

import (
	"fmt"
	"github.com/elliotcourant/noahdb/pkg/drivers/rqliter"
	"github.com/elliotcourant/noahdb/pkg/frunk"
	"github.com/elliotcourant/timber"
	"github.com/readystock/goqu"
)

var (
	getTenantsQuery = goqu.
		From("tenants").
		Select("tenants.*")
)

type tenantContext struct {
	*base
}

type TenantContext interface {
	GetTenants() ([]Tenant, error)
	GetTenant(uint64) (Tenant, error)
	NewTenants(tenantIds ...uint64) ([]Tenant, error)
}

func (ctx *base) Tenants() TenantContext {
	return &tenantContext{
		ctx,
	}
}

func (ctx *tenantContext) NewTenants(tenantIds ...uint64) ([]Tenant, error) {
	shardPressures, err := ctx.Shards().GetShardPressures(len(tenantIds))
	if err != nil {
		return nil, err
	}

	if len(shardPressures) == 0 {
		return nil, fmt.Errorf("no shards available for new tenants")
	}

	newTenants := make([]Tenant, len(tenantIds))
	for i, id := range tenantIds {
		shard := shardPressures[i%len(shardPressures)]
		newTenants[i] = Tenant{
			TenantID: id,
			ShardID:  shard.ShardID,
		}
		timber.Debugf("assigning new tenant [%d] to shard [%d]", id, shard.ShardID)
	}

	records := make([]interface{}, len(newTenants))
	for i, tenant := range newTenants {
		records[i] = goqu.Record{
			"tenant_id": tenant.TenantID,
			"shard_id":  tenant.ShardID,
		}
	}

	compiledSql := goqu.From("tenants").
		Insert(records...).Sql

	_, err = ctx.db.Exec(compiledSql)
	return newTenants, err
}

func (ctx *tenantContext) GetTenants() ([]Tenant, error) {
	compiledSql, _, _ := getTenantsQuery.
		ToSql()
	response, err := ctx.db.Query(compiledSql)
	if err != nil {
		return nil, err
	}
	return ctx.tenantsFromRows(response)
}

func (ctx *tenantContext) GetTenant(id uint64) (Tenant, error) {
	compiledSql, _, _ := getTenantsQuery.
		Where(goqu.Ex{
			"tenant_id": id,
		}).
		Limit(1).
		ToSql()
	response, err := ctx.db.Query(compiledSql)
	if err != nil {
		return Tenant{}, err
	}
	tenants, err := ctx.tenantsFromRows(response)
	if err != nil {
		return Tenant{}, err
	}
	if len(tenants) != 1 {
		return Tenant{}, fmt.Errorf("tried to find one tenant, found %d", len(tenants))
	}
	return tenants[0], nil
}

func (ctx *tenantContext) tenantsFromRows(response *frunk.QueryResponse) ([]Tenant, error) {
	rows := rqliter.NewRqlRows(response)
	tenants := make([]Tenant, 0)
	for rows.Next() {
		tenant := Tenant{}
		if err := rows.Scan(
			&tenant.TenantID,
			&tenant.ShardID); err != nil {
			return nil, err
		}
		tenants = append(tenants, tenant)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tenants, nil
}
