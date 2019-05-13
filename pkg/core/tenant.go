package core

import (
	"fmt"
	"github.com/elliotcourant/noahdb/pkg/drivers/rqliter"
	"github.com/elliotcourant/noahdb/pkg/frunk"
	"github.com/readystock/goqu"
)

var (
	getTenantsQuery = goqu.From("tenants").
		Select("tenants.*")
)

type tenantContext struct {
	*base
}

type TenantContext interface {
	GetTenants() ([]Tenant, error)
	GetTenant(uint64) (Tenant, error)
}

func (ctx *base) Tenants() TenantContext {
	return &tenantContext{
		ctx,
	}
}

func (ctx *tenantContext) GetTenants() ([]Tenant, error) {
	response, err := ctx.db.Query("SELECT tenant_id, shard_id FROM tenants;")
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
	return nil, nil
}
