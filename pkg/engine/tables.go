package engine

import (
	"errors"
	"github.com/elliotcourant/mellivora"
	"github.com/elliotcourant/noahdb/pkg/types"
)

var (
	// ErrInvalidTableName is returned when the table specified has no parts or has more than 2
	// parts.
	ErrInvalidTableName = errors.New("table name provided is not valid")

	// ErrTableDoesNotExist is returned when a table is requested by name or by schema and name and
	// no table matches the criteria.
	ErrTableDoesNotExist = errors.New("table does not exist")
)

var (
	_ TableContext = &tableContextBase{}
)

type (
	// TableType is used to indicate how NoahDb should handle queries directed to the table.
	TableType int

	// Table is a model that represents an actual table in PostgreSQL. Its used to keep track of
	// what tables have been created and what type they are to manage the database.
	Table struct {
		TableId     uint64 `m:"pk"`
		Schema      string `m:"uq:uq_schema_id_table_name"`
		Name        string `m:"uq:uq_schema_id_table_name"`
		Type        TableType
		HasSequence bool
	}

	// Column represents a PostgreSQL column, this model contains all of the data needed for NoahDb
	// to keep track of how data needs to be queried and returned.
	Column struct {
		ColumnId        uint64 `m:"pk"`
		TableId         uint64 `m:"uq:uq_table_id_column_name"`
		Type            types.Type
		Index           int
		Name            string `m:"uq:uq_table_id_column_name"`
		IsPrimaryKey    bool
		IsNullable      bool
		IsShardKey      bool
		IsSerial        bool
		ForeignColumnId uint64
	}

	// TableContext provides an accessor interface for tables within the cluster.
	TableContext interface {
		// NewTable will create a new table and the associated columns. If a table with the same
		// name already exists in the specified schema then an error will be returned.
		NewTable(table Table, columns []Column) (Table, []Column, error)

		// GetTableByName will return a table with the specified schema and name. If the schema is
		// not provided then it will assume "public". The name should be provided in parts. If you
		// only have the name of the table and not the schema then it will use that and sort by the
		// schema rank -> TODO (elliotcourant) add schema rank.
		// If you include the schema it should be called as GetTableByName("schema", "table").
		GetTableByName(name ...string) (Table, error)
	}

	tableContextBase struct {
		t *transactionBase
	}
)

const (
	// TableType_Unknown is used as a default value. Tables with this type should be ignored as they
	// are not properly setup.
	TableType_Unknown TableType = iota

	// TableType_Master indicates that the table is used to keep track of tenants within the
	// cluster. This table is typically an accounts or tenants table, the primary key of which will
	// be used to distribute and co-locate all of the data in the cluster.
	TableType_Master

	// TableType_Sharded indicates that the table's data is distributed in partitions throughout the
	// cluster. A single record in a sharded table will only exist within a single shard at any
	// point in time. That shard is then replicated as needed.
	TableType_Sharded

	// TableType_Global indicates that the table's data is present on every node and shard in the
	// entire cluster. The data in this table will be identical on every node and shard in the
	// cluster.
	TableType_Global

	// TableType_Postgres indicates that the table is not a user created table but is instead a
	// built-in table of the underlying PostgreSQL database.
	TableType_Postgres
)

// Tables will return the accessor interface for the table model.
func (t *transactionBase) Tables() TableContext {
	return &tableContextBase{
		t: t,
	}
}

// NewTable will create a new table and the associated columns. If a table with the same
// name already exists in the specified schema then an error will be returned.
func (t *tableContextBase) NewTable(table Table, columns []Column) (Table, []Column, error) {
	id, err := t.t.core.store.NextSequenceId("tables")
	if err != nil {
		return table, columns, err
	}
	table.TableId = id

	for i := range columns {
		columnId, err := t.t.core.store.NextSequenceId("columns")
		if err != nil {
			return table, columns, err
		}
		columns[i].TableId = id
		columns[i].ColumnId = columnId
	}

	if err := t.t.txn.Insert(table); err != nil {
		return table, columns, err
	}

	if err := t.t.txn.Insert(columns); err != nil {
		return table, columns, err
	}

	return table, columns, nil
}

// GetTableByName will return a table with the specified schema and name. If the schema is
// not provided then it will assume "public". The name should be provided in parts. If you
// only have the name of the table and not the schema then it will use that and sort by the
// schema rank -> TODO (elliotcourant) add schema rank.
// If you include the schema it should be called as GetTableByName("schema", "table").
func (t *tableContextBase) GetTableByName(name ...string) (Table, error) {
	getTableByNameQuery := t.t.txn.Model(Table{})

	// Handle the possible lengths for the name array.
	switch len(name) {
	case 2:
		// If two names were provided then we can assume that the first name is the schema and the
		// second name is the table name. So we can grab the first item here as the schema and then
		// fallthrough to grab the last item as the table name.
		getTableByNameQuery = getTableByNameQuery.AndWhere(mellivora.Ex{
			"Schema": name[0],
		})
		fallthrough
	case 1:
		// The last item of the provided names should be the table name for the query.
		getTableByNameQuery = getTableByNameQuery.AndWhere(mellivora.Ex{
			"Name": name[len(name)-1],
		})
	default:
		return Table{}, ErrInvalidTableName
	}

	// We want to scan the results into an array to properly assert whether or not there was more
	// than a single table that met the criteria or if there were no tables.
	tables := make([]Table, 0)

	// Run the query.
	if err := getTableByNameQuery.Select(&tables); err != nil {
		return Table{}, err
	}

	switch len(tables) {
	case 0:
		return Table{}, ErrTableDoesNotExist
	default:
		// TODO (elliotcourant) add schema prioritization here.
		fallthrough
	case 1:
		return tables[0], nil
	}
}
