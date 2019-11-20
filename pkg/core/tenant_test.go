package core_test

import (
	"github.com/ahmetb/go-linq/v3"
	"github.com/elliotcourant/noahdb/pkg/core"
	"github.com/elliotcourant/noahdb/testutils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTenantContext_NewTenants(t *testing.T) {
	colony, cleanup := testutils.NewPgTestColony(t)
	defer cleanup()

	shards, err := colony.Shards().GetShards()
	if !assert.NoError(t, err) {
		panic(err)
	}
	assert.NotEmpty(t, shards)

	numberOfTenants := 15

	newTenantsIds := make([]uint64, numberOfTenants)
	for i := 0; i < numberOfTenants; i++ {
		newTenantsIds[i] = uint64(i)
	}

	tenants, err := colony.Tenants().NewTenants(newTenantsIds...)
	if !assert.NoError(t, err) {
		panic(err)
	}

	assert.NotEmpty(t, tenants)
	assert.Equal(t, numberOfTenants, len(tenants))

	pressure, err := colony.Shards().GetShardPressures(0)
	if !assert.NoError(t, err) {
		panic(err)
	}
	assert.NotEmpty(t, pressure)
	assert.Equal(t, int64(numberOfTenants), linq.From(pressure).Select(func(i interface{}) interface{} {
		return i.(core.ShardPressure).Tenants
	}).SumInts())
}
