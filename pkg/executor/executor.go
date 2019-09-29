package executor

import (
	"github.com/elliotcourant/noahdb/pkg/core"
	"github.com/elliotcourant/noahdb/pkg/pgwirebase"
	"github.com/elliotcourant/noahdb/pkg/plan"
	"github.com/elliotcourant/timber"
	"time"
)

type Executor interface {
	Execute(plan plan.ExecutionPlan)
}

func NewExecutor(colony core.Colony) Executor {
	return &executorBase{
		colony: colony,
	}
}

type executorBase struct {
	transactionState TransactionState
	colony           core.Colony
	log              timber.Logger
}

type responsePipeline struct {
	conn core.PoolConnection
	err  error
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
		}(i, task)
	}
}
