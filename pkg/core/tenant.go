package core

import (
	"fmt"
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
	rows, err := ctx.db.Query("SELECT tenant_id, shard_id FROM tenants;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	tenants := make([]Tenant, 0)
	for rows.Next() {
		tenant := Tenant{}
		rows.Scan(&tenant.TenantID, &tenant.ShardID)
		tenants = append(tenants, tenant)
	}
	return tenants, nil
}

func (ctx *tenantContext) GetTenant(id uint64) (Tenant, error) {
	compiledSql, _, _ := getTenantsQuery.
		Where(goqu.Ex{
			"tenant_id": id,
		}).
		Limit(1).
		ToSql()
	rows, err := ctx.db.Query(compiledSql)
	if err != nil {
		return Tenant{}, err
	}
	defer rows.Close()
	tenants := make([]Tenant, 0)
	for rows.Next() {
		tenant := Tenant{}
		rows.Scan(&tenant.TenantID, &tenant.ShardID)
		tenants = append(tenants, tenant)
	}
	if len(tenants) != 1 {
		return Tenant{}, fmt.Errorf("tried to find one tenant, found %d", len(tenants))
	}
	return tenants[0], nil
}
