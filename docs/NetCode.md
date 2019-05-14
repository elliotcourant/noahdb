# Net Code

noahdb uses a single unified net code for all of it's operations.
This net code mimics PostgreSQL's own net code as closely as possible.
Meaning most client drivers in your language of choice should work
out of the box with noahdb. 

But noahdb's net code goes beyond just supporting existing protocols.
noahdb's raft implementation and internal rpc interface both use this
mimicked net code. When you start a noahdb coordinator, you will specify
a listener address. That address will be used for all communication
between that coordinator and others in the cluster. 

When a new connection is established, noahdb parses the first message
sent from the client and uses the Protocol Version from that message
to determine what type of connection it is. PostgreSQL, Raft or RPC.

See more here: https://www.postgresql.org/docs/current/protocol-message-formats.html

In the documentation above, the StartupMessage has a protocol version.
This is the value that is used to determine the connection type.

Currently noahdb supports protocol version 3.0 for all PostgreSQL connections.

# Flow

![flow](/docs/imgs/netflow.png "Network Flow")
