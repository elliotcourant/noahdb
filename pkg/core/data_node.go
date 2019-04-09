package core

import (
	"database/sql"
)

type dataNodeContext struct {
	*base
}

type DataNodeContext interface {
	GetDataNodes() ([]DataNode, error)
}

func (ctx *base) DataNodes() DataNodeContext {
	return &dataNodeContext{
		ctx,
	}
}

func (ctx *dataNodeContext) GetDataNodes() ([]DataNode, error) {
	rows, err := ctx.db.Query("SELECT data_node_id,address,port,healthy FROM data_nodes;")
	if err != nil {
		return nil, err
	}
	return ctx.dataNodesFromRows(rows)
}

func (ctx *dataNodeContext) dataNodesFromRows(rows *sql.Rows) ([]DataNode, error) {
	defer rows.Close()
	nodes := make([]DataNode, 0)
	for rows.Next() {
		node := DataNode{}
		if err := rows.Scan(&node.DataNodeID, &node.Address, &node.Port, &node.Healthy); err != nil {
			return nil, err
		}
		nodes = append(nodes, node)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return nodes, nil
}
