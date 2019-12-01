package engine_test

import (
	"github.com/ahmetb/go-linq/v3"
	"github.com/elliotcourant/noahdb/pkg/engine"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTenantBaseContext_NewTenant(t *testing.T) {
	t.Run("no shards", func(t *testing.T) {
		cluster, cleanup := NewTestCoreCluster(t, 1)
		defer cleanup()

		txn := cluster.Begin(t)

		_, err := txn.Tenants().NewTenant()
		assert.Equal(t, engine.ErrNoReadyShardsForTenant, err)
	})

	t.Run("simple", func(t *testing.T) {
		cluster, cleanup := NewTestCoreCluster(t, 1)
		defer cleanup()

		cluster.SeedDataNodes(t, 3)
		cluster.SeedShards(t, 3)

		txn := cluster.Begin(t)

		for i := 0; i < 3; i++ {
			tenant, err := txn.Tenants().NewTenant()
			assert.NoError(t, err)
			assert.NotZero(t, tenant.TenantId)
			assert.NotZero(t, tenant.ShardId)
		}
	})
}

func TestTenantBaseContext_GetTenant(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		cluster, cleanup := NewTestCoreCluster(t, 1)
		defer cleanup()

		cluster.SeedDataNodes(t, 1)
		cluster.SeedShards(t, 1)

		txn := cluster.Begin(t)

		tenant, err := txn.Tenants().NewTenant()
		assert.NoError(t, err)
		assert.NotZero(t, tenant.TenantId)
		assert.NotZero(t, tenant.ShardId)

		tenantRead, err := txn.Tenants().GetTenant(tenant.TenantId)
		assert.NoError(t, err)
		assert.Equal(t, tenant, tenantRead)
	})

	t.Run("not found", func(t *testing.T) {
		cluster, cleanup := NewTestCoreCluster(t, 1)
		defer cleanup()

		txn := cluster.Begin(t)

		_, err := txn.Tenants().GetTenant(1)
		assert.Equal(t, engine.ErrTenantNotFound, err)
	})
}

func TestTenantBaseContext_GetTenants(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		numberOfTenants := 3
		cluster, cleanup := NewTestCoreCluster(t, 1)
		defer cleanup()

		cluster.SeedDataNodes(t, 3)
		cluster.SeedShards(t, 3)

		txn := cluster.Begin(t)

		for i := 0; i < numberOfTenants; i++ {
			tenant, err := txn.Tenants().NewTenant()
			assert.NoError(t, err)
			assert.NotZero(t, tenant.TenantId)
			assert.NotZero(t, tenant.ShardId)
		}

		tenants, err := txn.Tenants().GetTenants()
		assert.NoError(t, err)
		assert.Len(t, tenants, numberOfTenants)
	})
}

func TestTenantBaseContext_GetTenantShardDistribution(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		numberOfShards, numberOfTenants := 3, 9
		cluster, cleanup := NewTestCoreCluster(t, 1)
		defer cleanup()

		cluster.SeedDataNodes(t, 3)
		cluster.SeedShards(t, numberOfShards)

		txn := cluster.Begin(t)

		for i := 0; i < numberOfTenants; i++ {
			tenant, err := txn.Tenants().NewTenant()
			assert.NoError(t, err)
			assert.NotZero(t, tenant.TenantId)
			assert.NotZero(t, tenant.ShardId)
		}

		tenants, err := txn.Tenants().GetTenantShardDistribution()
		assert.NoError(t, err)
		assert.Len(t, tenants, numberOfShards)
		assert.Equal(t,
			numberOfTenants,
			int(linq.From(tenants).
				Select(func(i interface{}) interface{} {
					return i.(linq.KeyValue).Value
				}).
				SumInts()),
		)
	})
}
