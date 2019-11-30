package engine

type (
	// Tenant represents a single group of co-located data. Queries must include a tenant Id filter
	// which will be used to determine what shard a given query should be directed to.
	Tenant struct {
		TenantId uint64 `m:"pk"`
		ShardId  uint64
	}

	TenantContext interface {
		// GetTenants will return all of the tenants in the entire cluster.
		GetTenants() ([]Tenant, error)

		// GetTenant will lookup a single specific tenant by their Id, if the tenant does not exist
		// then an ErrTenantNotFound error will be returned.
		GetTenant(tenantId uint64) (Tenant, error)

		// NewTenant will create a new tenant record and pick a shard from the cluster that the
		// tenant could be placed on.
		NewTenant(tenantId uint64) ([]Tenant, error)
	}
)
