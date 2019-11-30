package engine_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCoreBase_Begin(t *testing.T) {
	cluster, cleanup := NewTestCoreCluster(t, 3)
	defer cleanup()
	txn, err := cluster[0].Begin()
	assert.NoError(t, err)
	assert.NotNil(t, txn)
}
