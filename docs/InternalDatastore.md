# Internal Datastore

noahdb uses two different internal data-stores. A Badger key-value store
and a SQLite database. The key-value store handles record ID generation
and Raft information. The SQLite database is used as a structured mirror
of the key-value store. As actions occur in the cluster they are translated
into a protobuf message with the information about that action, they are
then distributed throughout the raft cluster. As actions are received on
each coordinator; if the action is a SQL action then the query will be
applied to the local SQLite database. The action is also stored in the
raft log and they key-value store. This means that when a new coordinator
comes online, or needs to restore from a snapshot: we replay the SQL queries
from the raft log onto the SQLite database. This is not terribly efficient
but is incredibly simple to maintain, and as long as all of the queries
directed to the local SQLite database go through our own abstraction layer
we can maintain the structure of the database.

## SQL Store

The SQL store is used to manage some of the basic topography data that noah
needs to know in order to manage the cluster. 

The definition for this internal database can be found here: 
[Internal Schema](/pkg/core/static/files/00_internal_sql.sql)

If this file is changed then the following command should be run.

```bash
make embedded
```

This will take the contents of the `files` folder that the SQL script resides
in and it will embed those files into a file named `static.embedded.go` (this
file is git-ignored) and the script is executed when the coordinator node
initially starts. 

If an internal store has already been established for a given coordinator then
the script will not be executed when the node starts. At the moment this is
intentional and is used for testing purposes. In the future there will be a
migration tool that will be used to migrate the internal database without causing
downtime within the cluster.

## Key-Value Store