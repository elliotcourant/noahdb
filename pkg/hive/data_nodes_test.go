package hive_test

import (
	"github.com/elliotcourant/meles"
	"github.com/elliotcourant/noahdb/pkg/core"
	"github.com/elliotcourant/noahdb/pkg/hive"
	"github.com/elliotcourant/noahdb/testutils"
	"github.com/elliotcourant/timber"
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
)

func NewHive(t *testing.T) (hive.Core, func()) {
	tempDir, cleanup := testutils.TempFolder()
	ln, err := net.Listen("tcp", ":")
	if !assert.NoError(t, err) {
		panic(err)
	}
	h, err := hive.NewHive(ln, timber.With(timber.Keys{
		"test": t.Name(),
	}), meles.Options{
		Directory: tempDir,
		Peers:     []string{},
	})
	if !assert.NoError(t, err) {
		panic(err)
	}
	return h, cleanup
}

func TestDataNodesContext_AddDataNode(t *testing.T) {
	h, cleanup := NewHive(t)
	defer cleanup()
	err := h.Start()
	assert.NoError(t, err)

	txn, err := h.Begin()
	assert.NoError(t, err)

	node, err := txn.DataNodes().AddDataNode(core.DataNode{
		Address:  "127.0.0.1",
		Port:     5432,
		User:     "postgres",
		Password: "password",
		Healthy:  false,
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, node)
}

func TestDataNodesContext_GetDataNodes(t *testing.T) {
	h, cleanup := NewHive(t)
	defer cleanup()
	err := h.Start()
	assert.NoError(t, err)

	txn, err := h.Begin()
	assert.NoError(t, err)

	node, err := txn.DataNodes().AddDataNode(core.DataNode{
		Address:  "127.0.0.1",
		Port:     5432,
		User:     "postgres",
		Password: "password",
		Healthy:  false,
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, node)

	node, err = txn.DataNodes().AddDataNode(core.DataNode{
		Address:  "127.0.0.1",
		Port:     5432,
		User:     "postgres",
		Password: "password",
		Healthy:  false,
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, node)

	nodes, err := txn.DataNodes().GetDataNodes()
	assert.NoError(t, err)
	assert.NotNil(t, nodes)
	assert.Len(t, nodes, 2)
}

func TestDataNodesContext_GetDataNode(t *testing.T) {
	h, cleanup := NewHive(t)
	defer cleanup()
	err := h.Start()
	assert.NoError(t, err)

	txn, err := h.Begin()
	assert.NoError(t, err)

	node, err := txn.DataNodes().AddDataNode(core.DataNode{
		Address:  "127.0.0.1",
		Port:     5432,
		User:     "postgres",
		Password: "password",
		Healthy:  false,
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, node)

	n, ok, err := txn.DataNodes().GetDataNode(node.DataNodeID)
	assert.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, node.DataNodeID, n.DataNodeID)
}
