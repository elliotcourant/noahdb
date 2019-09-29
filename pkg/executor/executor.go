package executor

import (
	"github.com/elliotcourant/noahdb/pkg/ast"
	"github.com/elliotcourant/noahdb/pkg/core"
	"github.com/elliotcourant/noahdb/pkg/pgproto"
	"github.com/elliotcourant/noahdb/pkg/pgwirebase"
	"github.com/elliotcourant/noahdb/pkg/plan"
	"github.com/elliotcourant/noahdb/pkg/pool"
	"github.com/elliotcourant/timber"
	"time"
)

type responsePipeline struct {
	conn core.PoolConnection
	err  error
}

type Executor interface {
	SetExtendedQueryMode(extended bool)
	Execute(plan plan.ExecutionPlan)
}

func NewExecutor(colony core.Colony, logger timber.Logger) Executor {
	return &executorBase{
		colony: colony,
		log:    logger,
		pool:   pool.NewPool(colony, logger),
	}
}

type executorBase struct {
	inTransaction     bool
	colony            core.Colony
	log               timber.Logger
	pool              pool.Pool
	extendedQueryMode bool
}

func (e *executorBase) Begin() {
	e.inTransaction = true
}

func (e *executorBase) SetExtendedQueryMode(extended bool) {
	e.extendedQueryMode = true
}

func (e *executorBase) Execute(executionPlan plan.ExecutionPlan) {
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

			conn, err := e.pool.GetConnection(task.DataNodeShardID)
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

			if extendedQuery {

			} else {
				if err := conn.Send(&pgproto.Query{
					String: task.Query,
				}); err != nil {
					e.log.Errorf("could not send query to data node shard [%d]: %v",
						task.DataNodeShardID,
						err)
					response.err = err
					return
				}
			}
		}(i, task)
	}
}
