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

The SQL store is used to help noah process some very simple queries when
a Postgres node is not available, or when a query is so simple that it's
easier to use the local SQLite database rather than sending a query over
the network for an additional round trip. It's primary purpose though is
to manage noah's cluster interface via SQL commands. After starting a
coordinator that is not connected to any Postgres nodes you can still
issue queries to add/remove data node, view information about cluster
topography.

### Tables

> #### noah.nodes
> Table definition:
> ```postgresql
> CREATE TABLE noah.data_nodes (
>   data_node_id SERIAL PRIMARY KEY,
>   address      TEXT   NOT NULL,
>   port         INT    NOT NULL,
>   healthy      BOOL   NOT NULL
> );
> ```

## Key-Value Store