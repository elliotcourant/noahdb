package engine

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
		ColumnId uint64 `m:"pk"`
		TableId  uint64 `m:"uq:uq_table_id_column_name"`
		// TODO (elliotcourant) Add Type field.
		Index           int
		Name            string `m:"uq:uq_table_id_column_name"`
		IsPrimaryKey    bool
		IsNullable      bool
		IsShardKey      bool
		IsSerial        bool
		ForeignColumnId uint64
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
