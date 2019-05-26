PRAGMA foreign_keys = ON;

CREATE TABLE users (
    user_id   BIGINT PRIMARY KEY,
    user_name TEXT NOT NULL UNIQUE,
    password  TEXT NOT NULL
);

CREATE TABLE types (
    type_id        INT PRIMARY KEY, -- type_id is the equivalent of postgres's's's's OID system.
    type_name      TEXT    NOT NULL UNIQUE,
    extension_type BOOLEAN NOT NULL
);

CREATE TABLE type_aliases (
    alias_id   INT PRIMARY KEY,
    type_id    INT  NOT NULL REFERENCES types (type_id) ON DELETE CASCADE,
    alias_name TEXT NOT NULL UNIQUE
);

CREATE TABLE settings (
    setting_id    BIGINT PRIMARY KEY,
    setting_key   TEXT    NOT NULL UNIQUE,
    type_id       BIGINT  NOT NULL REFERENCES types (type_id),
    int_value     BIGINT  NULL,
    boolean_value BOOLEAN NULL,
    text_value    TEXT    NULL
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
    primary_key BOOLEAN NOT NULL,
    nullable    BOOLEAN NOT NULL,
    shard_key   BOOLEAN NOT NULL,
    serial      BOOLEAN NOT NULL,
    UNIQUE (table_id, column_name)
);

INSERT INTO types (type_id, type_name, extension_type)
VALUES (16, 'bool', false),
       (17, 'bytea', false),
       (18, 'char', false),
       (19, 'name', false),
       (20, 'int8', false),
       (21, 'int2', false),
       (22, 'int2vector', false),
       (23, 'int4', false),
       (24, 'regproc', false),
       (25, 'text', false),
       (26, 'oid', false),
       (27, 'tid', false),
       (28, 'xid', false),
       (29, 'cid', false),
       (30, 'oidvector', false),
       (114, 'json', false),
       (142, 'xml', false),
       (143, '_xml', false),
       (199, '_json', false),
       (194, 'pg_node_tree', false),
       (3361, 'pg_ndistinct', false),
       (3402, 'pg_dependencies', false),
       (32, 'pg_ddl_command', false),
       (210, 'smgr', false),
       (600, 'point', false),
       (601, 'lseg', false),
       (602, 'path', false),
       (603, 'box', false),
       (604, 'polygon', false),
       (628, 'line', false),
       (629, '_line', false),
       (700, 'float4', false),
       (701, 'float8', false),
       (702, 'abstime', false),
       (703, 'reltime', false),
       (704, 'tinterval', false),
       (705, 'unknown', false),
       (718, 'circle', false),
       (719, '_circle', false),
       (790, 'money', false),
       (791, '_money', false),
       (829, 'macaddr', false),
       (869, 'inet', false),
       (650, 'cidr', false),
       (774, 'macaddr8', false),
       (1000, '_bool', false),
       (1001, '_bytea', false),
       (1002, '_char', false),
       (1003, '_name', false),
       (1005, '_int2', false),
       (1006, '_int2vector', false),
       (1007, '_int4', false),
       (1008, '_regproc', false),
       (1009, '_text', false),
       (1028, '_oid', false),
       (1010, '_tid', false),
       (1011, '_xid', false),
       (1012, '_cid', false),
       (1013, '_oidvector', false),
       (1014, '_bpchar', false),
       (1015, '_varchar', false),
       (1016, '_int8', false),
       (1017, '_point', false),
       (1018, '_lseg', false),
       (1019, '_path', false),
       (1020, '_box', false),
       (1021, '_float4', false),
       (1022, '_float8', false),
       (1023, '_abstime', false),
       (1024, '_reltime', false),
       (1025, '_tinterval', false),
       (1027, '_polygon', false),
       (1033, 'aclitem', false),
       (1034, '_aclitem', false),
       (1040, '_macaddr', false),
       (775, '_macaddr8', false),
       (1041, '_inet', false),
       (651, '_cidr', false),
       (1263, '_cstring', false),
       (1042, 'bpchar', false),
       (1043, 'varchar', false),
       (1082, 'date', false),
       (1083, 'time', false),
       (1114, 'timestamp', false),
       (1115, '_timestamp', false),
       (1182, '_date', false),
       (1183, '_time', false),
       (1184, 'timestamptz', false),
       (1185, '_timestamptz', false),
       (1186, 'interval', false),
       (1187, '_interval', false),
       (1231, '_numeric', false),
       (1266, 'timetz', false),
       (1270, '_timetz', false),
       (1560, 'bit', false),
       (1561, '_bit', false),
       (1562, 'varbit', false),
       (1563, '_varbit', false),
       (1700, 'numeric', false),
       (1790, 'refcursor', false),
       (2201, '_refcursor', false),
       (2202, 'regprocedure', false),
       (2203, 'regoper', false),
       (2204, 'regoperator', false),
       (2205, 'regclass', false),
       (2206, 'regtype', false),
       (4096, 'regrole', false),
       (4089, 'regnamespace', false),
       (2207, '_regprocedure', false),
       (2208, '_regoper', false),
       (2209, '_regoperator', false),
       (2210, '_regclass', false),
       (2211, '_regtype', false),
       (4097, '_regrole', false),
       (4090, '_regnamespace', false),
       (2950, 'uuid', false),
       (2951, '_uuid', false),
       (3220, 'pg_lsn', false),
       (3221, '_pg_lsn', false),
       (3614, 'tsvector', false),
       (3642, 'gtsvector', false),
       (3615, 'tsquery', false),
       (3734, 'regconfig', false),
       (3769, 'regdictionary', false),
       (3643, '_tsvector', false),
       (3644, '_gtsvector', false),
       (3645, '_tsquery', false),
       (3735, '_regconfig', false),
       (3770, '_regdictionary', false),
       (3802, 'jsonb', false),
       (3807, '_jsonb', false),
       (2970, 'txid_snapshot', false),
       (2949, '_txid_snapshot', false),
       (3904, 'int4range', false),
       (3905, '_int4range', false),
       (3906, 'numrange', false),
       (3907, '_numrange', false),
       (3908, 'tsrange', false),
       (3909, '_tsrange', false),
       (3910, 'tstzrange', false),
       (3911, '_tstzrange', false),
       (3912, 'daterange', false),
       (3913, '_daterange', false),
       (3926, 'int8range', false),
       (3927, '_int8range', false),
       (2249, 'record', false),
       (2287, '_record', false),
       (2275, 'cstring', false),
       (2276, 'any', false),
       (2277, 'anyarray', false),
       (2278, 'void', false),
       (2279, 'trigger', false),
       (3838, 'event_trigger', false),
       (2280, 'language_handler', false),
       (2281, 'internal', false),
       (2282, 'opaque', false),
       (2283, 'anyelement', false),
       (2776, 'anynonarray', false),
       (3500, 'anyenum', false),
       (3115, 'fdw_handler', false),
       (325, 'index_am_handler', false),
       (3310, 'tsm_handler', false),
       (3831, 'anyrange', false),
       (141313, 'hstore', false),
       (141318, '_hstore', false),
       (141397, 'ghstore', false),
       (141400, '_ghstore', false);

INSERT INTO type_aliases (alias_id, type_id, alias_name)
VALUES (1, 16, 'boolean'),
       (2, 20, 'bigint'),
       (3, 1016, '_bigint'),
       (4, 21, 'smallint'),
       (5, 1005, '_smallint'),
       (6, 23, 'int'),
       (7, 1007, '_int'),
       (8, 23, 'integer'),
       (9, 1007, '_integer'),
       (10, 25, 'string'),
       (11, 1184, 'timestamp with time zone'),
       (12, 1185, '_timestamp with time zone'),
       (13, 1266, 'time with time zone'),
       (14, 1270, '_time with time zone');

INSERT INTO settings (setting_id, setting_key, type_id, int_value, boolean_value, text_value)
VALUES (1, 'replication_mode', 20, 0, null, null),
       (2, 'replication_factor', 20, 1, null, null),
       (3, 'max_pool_size', 20, 5, null, null),
       (4, 'min_pool_size', 20, 0, null, null),
       (5, 'pool_refresh_interval', 1186, null, null, '30 seconds'),
       (6, 'number_of_shards', 20, 3, null, null);