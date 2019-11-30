package engine_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCoreBase_Begin(t *testing.T) {
	cluster, cleanup := NewTestCoreCluster(t, 1)
	defer cleanup()
	txn, err := cluster[0].Begin()
	assert.NoError(t, err)
	assert.NotNil(t, txn)
}

func TestCoreBase_Rollback(t *testing.T) {
	cluster, cleanup := NewTestCoreCluster(t, 1)
	defer cleanup()
	txn, err := cluster[0].Begin()
	assert.NoError(t, err)
	assert.NotNil(t, txn)

	_, err = txn.
		DataNodes().
		NewDataNode("test", 1234, "test", "test")
	assert.NoError(t, err)

	err = txn.Rollback()
	assert.NoError(t, err)

	// If we try to use the transaction now that it has been disposed it should return an error.
	_, err = txn.
		DataNodes().
		NewDataNode("test", 1234, "test", "test")
	assert.Error(t, err)
}
