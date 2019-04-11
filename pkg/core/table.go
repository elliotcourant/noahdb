package core

type tableContext struct {
	*base
}

type TableContext interface {
	GetTables() ([]Table, error)
}

func (ctx *base) Tables() TableContext {
	return &tableContext{
		ctx,
	}
}

func (ctx *tableContext) GetTables() ([]Table, error) {
	return nil, nil
}

func (ctx *tableContext) GetTable(schemaId int, tableName string) (Table, error) {
	return Table{}, nil
}
