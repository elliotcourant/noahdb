PRAGMA foreign_keys = ON;

CREATE TABLE users (
    user_id   BIGINT PRIMARY KEY,
    user_name TEXT NOT NULL UNIQUE,
    password  TEXT NOT NULL
);

CREATE TABLE types (
    type_id        INT PRIMARY KEY, -- type_id is the equivalent of postgres's's's's OID system.
    type_name      TEXT    NOT NULL UNIQUE,
    postgres_name  TEXT    NOT NULL,
    extension_type BOOLEAN NOT NULL,
    alias_type_id  INT     NULL REFERENCES types (type_id) ON DELETE CASCADE
);

CREATE TABLE settings (
    setting_id BIGINT PRIMARY KEY,
    type_id    BIGINT NOT NULL REFERENCES types (type_id)
);

CREATE TABLE data_nodes (
    data_node_id BIGINT PRIMARY KEY,
    address      TEXT    NOT NULL,
    port         INT     NOT NULL,
    healthy      BOOLEAN NOT NULL
);

CREATE TABLE shards (
    shard_id BIGINT PRIMARY KEY,
    state    INT NOT NULL
);

CREATE TABLE tenants (
    tenant_id BIGINT PRIMARY KEY,
    shard_id  BIGINT NOT NULL,
    FOREIGN KEY (shard_id) REFERENCES shards (shard_id)
);

CREATE TABLE data_node_shards (
    data_node_shard_id BIGINT PRIMARY KEY,
    data_node_id       BIGINT  NOT NULL,
    shard_id           BIGINT  NOT NULL,
    read_only          BOOLEAN NOT NULL,
    UNIQUE (data_node_id, shard_id),
    FOREIGN KEY (data_node_id) REFERENCES data_nodes (data_node_id),
    FOREIGN KEY (shard_id) REFERENCES shards (shard_id)
);



CREATE TABLE schemas (
    schema_id   INT PRIMARY KEY,
    schema_name TEXT NOT NULL UNIQUE
);

CREATE TABLE tables (
    table_id   INT PRIMARY KEY,
    schema_id  INT  NOT NULL REFERENCES schemas (schema_id) ON DELETE CASCADE,
    table_name TEXT NOT NULL,
    table_type INT  NOT NULL,
    UNIQUE (schema_id, table_name)
);

CREATE TABLE columns (
    column_id   BIGINT PRIMARY KEY,
    table_id    INT     NOT NULL REFERENCES tables (table_id) ON DELETE CASCADE,
    type_id     INT     NOT NULL REFERENCES types (type_id) ON DELETE RESTRICT,
    sort        INT     NOT NULL,
    column_name TEXT    NOT NULL,
    nullable    BOOLEAN NOT NULL,
    shard_key   BOOLEAN NOT NULL,
    serial      BOOLEAN NOT NULL,
    UNIQUE (table_id, column_name)
);

