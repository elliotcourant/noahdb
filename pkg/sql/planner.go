package sql

import (
	"fmt"
	"github.com/elliotcourant/noahdb/pkg/ast"
	"github.com/elliotcourant/noahdb/pkg/pgwirebase"
	"time"
)

type PlanType string

const (
	PlanType_READ      PlanType = "READ"
	PlanType_WRITE     PlanType = "WRITE"
	PlanType_READWRITE PlanType = "READWRITE"
)

type PlanTarget string

const (
	PlanTarget_STANDARD PlanTarget = "STANDARD"
	PlanTarget_INTERNAL PlanTarget = "INTERNAL"
)

type DistributedPlanType int

const (
	DistributedPlanType_NONE     DistributedPlanType = 0
	DistributedPlanType_COMMIT   DistributedPlanType = 1
	DistributedPlanType_ROLLBACK DistributedPlanType = 2
)

type InitialPlanTask struct {
	Query string
	Type  ast.StmtType
}

type InitialPlan struct {
	Types        map[PlanType]InitialPlanTask
	ShardID      uint64
	Target       PlanTarget
	DistPlanType DistributedPlanType
}

type ExpandedPlan struct {
	Tasks        []ExpandedPlanTask
	Target       PlanTarget
	OutFormats   []pgwirebase.FormatCode
	DistPlanType DistributedPlanType
}

type ExpandedPlanTask struct {
	Query           string
	ReadOnly        bool
	DataNodeShardID uint64
	Type            ast.StmtType
}

type NoahQueryPlanner interface {
	getNoahQueryPlan(s *session) (InitialPlan, bool, error)
}

type QueryPlanner interface {
	GetQueryPlan(s *session) (InitialPlan, bool, error)
}

type TransactionQueryPlanner interface {
	getTransactionQueryPlan(s *session) (InitialPlan, bool, error)
}

func (s *session) expandQueryPlan(plan InitialPlan) (ExpandedPlan, error) {
	startTimestamp := time.Now()
	defer func() {
		s.log.Verbosef("[%s] expanding of plan", time.Since(startTimestamp))
	}()

	if plan.DistPlanType != DistributedPlanType_NONE {
		pendingDataNodeShards := s.GetPendingDataNodeShards()
		// useTwoPhaseCommit := len(pendingDataNodeShards) > 1
		tasks := make([]ExpandedPlanTask, len(pendingDataNodeShards))

		for i, id := range pendingDataNodeShards {
			switch plan.DistPlanType {
			case DistributedPlanType_COMMIT:
				tasks[i] = ExpandedPlanTask{
					Query:           "COMMIT",
					ReadOnly:        false,
					DataNodeShardID: id,
					Type:            ast.Ack,
				}
			case DistributedPlanType_ROLLBACK:
				tasks[i] = ExpandedPlanTask{
					Query:           "ROLLBACK",
					ReadOnly:        false,
					DataNodeShardID: id,
					Type:            ast.Ack,
				}
			}
		}

		return ExpandedPlan{
			Target:       PlanTarget_STANDARD,
			Tasks:        tasks,
			DistPlanType: plan.DistPlanType,
		}, nil
	}

	if plan.Target == PlanTarget_INTERNAL {
		// Internal query plans can go directly to the SQLite database.
		s.log.Verbosef("plan targets internal SQLite database")
		readPlan, _ := plan.Types[PlanType_READ]
		return ExpandedPlan{
			Target: PlanTarget_INTERNAL,
			Tasks: []ExpandedPlanTask{
				{
					Query:           readPlan.Query,
					ReadOnly:        true,
					DataNodeShardID: 0,
					Type:            readPlan.Type,
				},
			},
		}, nil
	}

	readOnly := true
	dataNodeShards := make([]uint64, 0)
	switch plan.ShardID {
	case 0: // If this query does not target a specific shard.
		// If we are performing a read then the planner should have accommodated for the read being
		// sent to any shard, otherwise a shard would be specified. So we can send this query to
		// any random node.
		if _, ok := plan.Types[PlanType_READ]; ok {
			// Get a single node to execute the read query.
			id, err := s.Colony().DataNodes().GetRandomDataNodeShardID()
			if err != nil {
				return ExpandedPlan{}, err
			}
			dataNodeShards = append(dataNodeShards, id)
			break
		}

		// If we are performing a write, and there is no shard ID then that means the write must be
		// performed on ALL shards in the cluster.
		if _, ok := plan.Types[PlanType_WRITE]; ok {
			readOnly = false
		} else if _, ok := plan.Types[PlanType_READWRITE]; ok {
			readOnly = false
		}

		if !readOnly {
			ids, err := s.Colony().DataNodes().GetDataNodeShardIDs()
			if err != nil {
				return ExpandedPlan{}, err
			}
			dataNodeShards = append(dataNodeShards, ids...)
			break
		}
	default:
		tempDataNodeShardIds, err := s.Colony().DataNodes().GetDataNodeShardIDsForShard(plan.ShardID)
		if err != nil {
			return ExpandedPlan{}, fmt.Errorf("could not retrieve data nodes for shard ID [%d]: %s", plan.ShardID, err.Error())
		}

		if len(tempDataNodeShardIds) < 1 {
			return ExpandedPlan{}, fmt.Errorf("could not retrieve data nodes for shard ID [%d]: no nodes were returned", plan.ShardID)
		}

		if _, ok := plan.Types[PlanType_READ]; ok {
			// Get any readable node for the given shard ID
			dataNodeShards = append(dataNodeShards, tempDataNodeShardIds[0])
			break
		}

		if _, ok := plan.Types[PlanType_WRITE]; ok {
			readOnly = false
		} else if _, ok := plan.Types[PlanType_READWRITE]; ok {
			readOnly = false
		}

		if !readOnly {
			// Get all the nodes that are writable for the given shard ID
			dataNodeShards = tempDataNodeShardIds
			break
		}
	}

	// For each node/shard we are targeting, generate a task for the executor
	tasks := make([]ExpandedPlanTask, len(dataNodeShards))
	for i, id := range dataNodeShards {
		task := ExpandedPlanTask{
			ReadOnly:        readOnly,
			DataNodeShardID: id,
		}

		if readPlan, ok := plan.Types[PlanType_READ]; ok {
			task.Query, task.Type = readPlan.Query, readPlan.Type
		} else if writePlan, ok := plan.Types[PlanType_WRITE]; ok {
			task.Query, task.Type = writePlan.Query, writePlan.Type
		} else if writePlan, ok := plan.Types[PlanType_READWRITE]; ok {
			task.Query, task.Type = writePlan.Query, writePlan.Type
		}

		tasks[i] = task
	}

	if len(tasks) < 1 {
		return ExpandedPlan{}, fmt.Errorf("could not generate tasks for plan")
	}

	// If there is a returning clause on any of the query plan types
	// then we want to make sure that just one of the executions
	// has a returning clause.
	if readWritePlan, ok := plan.Types[PlanType_READWRITE]; ok {
		tasks[0].Query, tasks[0].Type = readWritePlan.Query, readWritePlan.Type
	}

	return ExpandedPlan{
		Target: plan.Target,
		Tasks:  tasks,
	}, nil
}
