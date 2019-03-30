package core

type tenantContext struct {
	*base
}

type TenantContext interface {
	GetTenants() ([]Tenant, error)
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
