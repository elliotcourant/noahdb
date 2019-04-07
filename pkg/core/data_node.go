package core

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
	defer rows.Close()
	nodes := make([]DataNode, 0)
	for rows.Next() {
		node := DataNode{}
		rows.Scan(&node.DataNodeID, &node.Address, &node.Port, &node.Healthy)
		nodes = append(nodes, node)
	}
	return nodes, nil
}
