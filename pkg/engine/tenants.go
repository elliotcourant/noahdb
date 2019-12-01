package engine

import (
	"errors"
	"github.com/ahmetb/go-linq/v3"
	"github.com/elliotcourant/mellivora"
)

var (
	// ErrTenantNotFound is returned when a tenant is requested specifically by it's Id
	// but a record with that Id does not exist.
	ErrTenantNotFound = errors.New("tenant does not exist")

	// ErrNoReadyShardsForTenant is returned when a tenant is being created but there are no
	// shards in the cluster that have a state of ShardState_Ready.
	ErrNoReadyShardsForTenant = errors.New("no ready shards to place tenant on")
)

var (
	_ TenantContext = &tenantBaseContext{}
)

type (
	// Tenant represents a single group of co-located data. Queries must include a tenant Id filter
	// which will be used to determine what shard a given query should be directed to.
	Tenant struct {
		TenantId uint64 `m:"pk"`
		ShardId  uint64
	}

	// TenantContext provides an accessor interface for tenants within the cluster.
	TenantContext interface {
		// GetTenants will return all of the tenants in the entire cluster.
		GetTenants() ([]Tenant, error)

		// GetTenant will lookup a single specific tenant by their Id, if the tenant does not exist
		// then an ErrTenantNotFound error will be returned.
		GetTenant(tenantId uint64) (Tenant, error)

		// GetTenantShardDistribution returns a map of shardIds and the number of tenants on each shard.
		GetTenantShardDistribution() (map[uint64]int, error)

		// NewTenant will create a new tenant record and pick a shard from the cluster that the
		// tenant could be placed on.
		NewTenant() (Tenant, error)
	}

	tenantBaseContext struct {
		t *transactionBase
	}
)

// Tenants returns the accessors for tenant models.
func (t *transactionBase) Tenants() TenantContext {
	return &tenantBaseContext{
		t: t,
	}
}

// GetTenants will return all of the tenants in the entire cluster.
func (t *tenantBaseContext) GetTenants() ([]Tenant, error) {
	tenants := make([]Tenant, 0)
	err := t.t.txn.Model(tenants).Select(&tenants)

	return tenants, err
}

// GetTenant will lookup a single specific tenant by their Id, if the tenant does not exist
// then an ErrTenantNotFound error will be returned.
func (t *tenantBaseContext) GetTenant(tenantId uint64) (Tenant, error) {
	tenant := Tenant{}
	err := t.t.txn.
		Model(tenant).
		Where(mellivora.Ex{
			"TenantId": tenantId,
		}).Select(&tenant)
	if tenant.TenantId == 0 && err == nil {
		return tenant, ErrTenantNotFound
	}

	return tenant, err
}

// NewTenant will create a new tenant record and pick a shard from the cluster that the
// tenant could be placed on.
func (t *tenantBaseContext) NewTenant() (Tenant, error) {
	tenant := Tenant{}

	distribution, err := t.GetTenantShardDistribution()
	if err != nil {
		return tenant, err
	} else if len(distribution) == 0 {
		return tenant, ErrNoReadyShardsForTenant
	}

	id, err := t.t.core.store.NextSequenceId("tenants")
	if err != nil {
		return tenant, err
	}

	tenant.TenantId = id

	// Find the shardId that this tenant should be placed on.
	tenant.ShardId = linq.From(distribution).
		// Sort the current shards by the number of tenants on each
		// one in ascending order.
		OrderBy(func(i interface{}) interface{} {
			return i.(linq.KeyValue).Value
		}).
		// Pull the shardId out.
		Select(func(i interface{}) interface{} {
			return i.(linq.KeyValue).Key
		}).
		// Grab the shard with the least number of tenants.
		First().(uint64)

	return tenant, t.t.txn.Insert(tenant)
}

// GetTenantShardDistribution returns a map of shardIds and the number of tenants on each shard.
func (t *tenantBaseContext) GetTenantShardDistribution() (map[uint64]int, error) {
	// Get ready shards that we can place a tenant on.
	shards, err := t.t.Shards().GetShards(ShardState_Ready)
	if err != nil {
		return nil, err
	}

	tenants, err := t.GetTenants()
	if err != nil {
		return nil, err
	}

	// Seed the distribution map with all the shards..
	distribution := map[uint64]int{}
	for _, shard := range shards {
		distribution[shard.ShardId] = 0
	}

	// For each tenant, look at what shard they belong to and increment
	// the number of tenants for that shard.
	for _, tenant := range tenants {
		distribution[tenant.ShardId] = distribution[tenant.ShardId] + 1
	}

	return distribution, nil
}
