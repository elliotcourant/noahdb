# noahdb

noahdb is a distributed multi-tenant new-sql database.

- General
    - [Requirements](/docs/README.md#requirements)
    - Installation
    - Build From Source
    - Configuration
    - Examples
- Documentation
    - Coordinators
    - Data Nodes
        - Workers
        - Full-Replicas
        - Analytics
    - Shards
    - Tenants
- Query Interface
    - Adding A Coordinator
    - Adding A Data Node
    - Adding A Shard
    - Get Tenants
        - Tenants On Each Shard
        - Tenants On Each Node
- Internals
    - [Parser](/docs/Parser.md)
    - [Internal Data-Store](/docs/InternalDatastore.md)
    
# Requirements

Wireshark is not a required tool, but is highly recommended when developing noahdb.
Protoc however, is a required tool if you want to do a clean build of noahdb.

- Wireshark (https://www.wireshark.org/)
    - Download: https://www.wireshark.org/#download
    
    
- Protoc (https://github.com/protocolbuffers/protobuf)
    - Download: http://google.github.io/proto-lens/installing-protoc.html
    
At the moment noahdb cannot be compiled on Windows, it must be compiled on Mac OS X or Linux.

# Commands
