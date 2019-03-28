package core

type shardContext struct {
	*base
}

type ShardContext interface {
	GetShards() ([]Shard, error)
}

func (ctx *base) Shards() ShardContext {
	return &shardContext{
		ctx,
	}
}

func (ctx *shardContext) GetShards() ([]Shard, error) {
	rows, err := ctx.db.Query("SELECT id FROM shards;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	shards := make([]Shard, 0)
	for rows.Next() {
		shard := Shard{}
		rows.Scan(&shard.ShardID)
		shards = append(shards, shard)
	}
	return shards, nil
}
