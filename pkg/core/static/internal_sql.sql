PRAGMA foreign_keys = ON;

CREATE TABLE data_nodes (
  data_node_id BIGINT PRIMARY KEY,
  address      TEXT    NOT NULL,
  port         INT     NOT NULL,
  healthy      BOOLEAN NOT NULL
);

CREATE TABLE shards (
  shard_id BIGINT PRIMARY KEY
);

CREATE TABLE tenants (
  tenant_id BIGINT PRIMARY KEY,
  shard_id  BIGINT NOT NULL,
  FOREIGN KEY (shard_id) REFERENCES shards (shard_id)
);

CREATE TABLE data_node_shards (
  data_node_shard_id BIGINT PRIMARY KEY,
  data_node_id       BIGINT NOT NULL,
  shard_id           BIGINT NOT NULL,
  FOREIGN KEY (data_node_id) REFERENCES data_nodes (data_node_id),
  FOREIGN KEY (shard_id) REFERENCES shards (shard_id)
);

