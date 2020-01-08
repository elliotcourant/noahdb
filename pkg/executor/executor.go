package executor

import (
	"github.com/elliotcourant/noahdb/pkg/ast"
	"github.com/elliotcourant/noahdb/pkg/engine"
	"github.com/elliotcourant/noahdb/pkg/pgproto"
	"github.com/elliotcourant/noahdb/pkg/pgwirebase"
	"github.com/elliotcourant/noahdb/pkg/plan"
	"github.com/elliotcourant/timber"
	"time"
)

type responsePipeline struct {
	conn engine.Connection
	err  error
}

type Executor interface {
	SetExtendedQueryMode(extended bool)
	Execute(plan plan.ExecutionPlan) error
}

func NewExecutor(txn engine.Transaction, logger timber.Logger, tunnel Tunnel) Executor {
	return &executorBase{
		txn:    txn,
		tunnel: tunnel,
		log:    logger,
		pool:   pool_old.NewPool(txn, logger),
	}
}

type executorBase struct {
	inTransaction     bool
	txn               engine.Transaction
	tunnel            Tunnel
	log               timber.Logger
	pool              pool_old.Pool
	extendedQueryMode bool
}

func (e *executorBase) Begin() {
	e.inTransaction = true
	e.pool.Begin() // Make sure to move the pool to be in a transaction
}

func (e *executorBase) CommitAndRelease() error {
	return nil
}

func (e *executorBase) SetExtendedQueryMode(extended bool) {
	e.extendedQueryMode = true
}

// Execute runs a plan against the data node cluster.
func (e *executorBase) Execute(executionPlan plan.ExecutionPlan) error {
	if len(executionPlan.OutFormats) == 0 {
		executionPlan.OutFormats = []pgwirebase.FormatCode{
			pgwirebase.FormatText,
		}
	}
	startTimestamp := time.Now()
	defer func() {
		e.log.Verbosef("[%s] execution of statement", time.Since(startTimestamp))
	}()

	responses := make(chan *responsePipeline, len(executionPlan.Tasks))

	for i, task := range executionPlan.Tasks {
		go func(index int, task plan.Task) {
			var response = &responsePipeline{}
			defer func() {
				e.log.Verbosef("[%s] dispatch of query to data node shard [%d]", time.Since(startTimestamp), task.DataNodeShardID)
				responses <- response
			}()

			// Retrieve a connection from the pool, if we are in a transaction or if there
			// are queries being executed against multiple shards then we need to start
			// a transaction.
			conn, err := e.pool.GetConnection(task.DataNodeShardID, len(executionPlan.Tasks) > 0)
			if err != nil {
				e.log.Errorf("could not retrieve connection from pool for data node shard [%d]: %v",
					task.DataNodeShardID, err)
				response.err = err
				return
			}
			e.log.Verbosef("{%d} executing: %s", task.DataNodeShardID, task.Query)

			extendedQuery := e.extendedQueryMode
			if task.Type != ast.Rows {
				extendedQuery = false
			}

			// Send the query to the target data node shard.
			response.err = e.sendQuery(conn, task.Query, extendedQuery, executionPlan.OutFormats)
			response.conn = conn
		}(i, task)
	}

	for i := 0; i < len(executionPlan.Tasks); i++ {
		if err := func(response *responsePipeline) error {
			if response.err != nil {
				return response.err
			}
			conn := response.conn
			if !e.inTransaction {
				// If we are not in a transaction, then we can yeet this connection back to
				// the global pool as soon as we process it's response.
				// TODO (elliotcourant) add pool return.
			}

			canExit := false
			for {
				message, err := conn.Receive()
				if err != nil {
					e.log.Errorf("received error from frontend: %v", err)
					return err
				}

				switch msg := message.(type) {
				case *pgproto.RowDescription, *pgproto.DataRow:
					if err := e.tunnel.SendMessage(msg); err != nil {
						return err
					}
					canExit = true
				case *pgproto.ErrorResponse:
					if err := e.tunnel.SendMessage(msg); err != nil {
						return err
					}
					canExit = true
					// Even though we received an error from the data node shard, this is not our
					// problem. At this point an error received would be due to a client mistake.
					// So we return nil here because this is not an error we can handle.
					// TODO (elliotcourant) add handling for inner-transaction errors.
					return nil
				case *pgproto.CommandComplete:
					canExit = true
					return nil
				case *pgproto.ReadyForQuery:
					if canExit {
						return nil
					}
				default:
					e.log.Warningf("received unexpected message [%T] from data node shard", msg)
				}
			}
		}(<-responses); err != nil {
			return err
		}
	}

	return nil
}

func (e *executorBase) sendQuery(conn engine.Connection, query string, extendedQuery bool, outFormats []pgwirebase.FormatCode) error {
	if extendedQuery {
		// When we are in extended query mode we want to send the query in the same
		// extended query mode.
		if err := conn.Send(&pgproto.Parse{
			Name:  "",
			Query: query,
		}); err != nil {
			e.log.Errorf(
				"could not send query to data node shard [%d]: %s",
				conn.DataNodeShardID(), err.Error())
			return err
		}

		if err := conn.Send(&pgproto.Describe{
			ObjectType: 'S',
			Name:       "",
		}); err != nil {
			e.log.Errorf(
				"could not describe query on data node shard [%d]: %s",
				conn.DataNodeShardID(), err.Error())
			return err
		}

		if err := conn.Send(&pgproto.Bind{
			DestinationPortal: "",
			PreparedStatement: "",
			ResultFormatCodes: outFormats,
		}); err != nil {
			e.log.Errorf(
				"could not bind on data node shard [%d]: %s",
				conn.DataNodeShardID(), err.Error())
			return err
		}

		if err := conn.Send(&pgproto.Execute{
			Portal:  "",
			MaxRows: 0,
		}); err != nil {
			e.log.Errorf(
				"could not bind on data node shard [%d]: %s",
				conn.DataNodeShardID(), err.Error())
			return err
		}

		if err := conn.Send(&pgproto.Sync{}); err != nil {
			e.log.Errorf(
				"could not sync on data node shard [%d]: %s",
				conn.DataNodeShardID(), err.Error())
			return err
		}
	} else {
		if err := conn.Send(&pgproto.Query{
			String: query,
		}); err != nil {
			e.log.Errorf("could not send query to data node shard [%d]: %v",
				conn.DataNodeShardID(), err)
			return err
		}
	}
	return nil
}
