package sql

import (
	"fmt"
	"github.com/elliotcourant/noahdb/pkg/ast"
	"github.com/readystock/golog"
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

type InitialPlanTask struct {
	Query string
	Type  ast.StmtType
}

type InitialPlan struct {
	Types   map[PlanType]InitialPlanTask
	ShardID uint64
	Target  PlanTarget
}

type ExpandedPlan struct {
	Tasks  []ExpandedPlanTask
	Target PlanTarget
}

type ExpandedPlanTask struct {
	Query           string
	ReadOnly        bool
	DataNodeShardID uint64
	Type            ast.StmtType
}

func (s *session) expandQueryPlan(plan InitialPlan) (ExpandedPlan, error) {
	startTimestamp := time.Now()
	defer func() {
		golog.Verbosef("expanding of plan took %s", time.Since(startTimestamp))
	}()

	if plan.Target == PlanTarget_INTERNAL {
		// Internal query plans can go directly to the SQLite database.
		golog.Verbosef("plan targets internal SQLite database")
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
		if _, ok := plan.Types[PlanType_READ]; ok {
			// Get a single node to execute the read query.
			id, err := s.Colony().DataNodes().GetRandomDataNodeShardID()
			if err != nil {
				return ExpandedPlan{}, err
			}
			dataNodeShards = append(dataNodeShards, id)
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
		} else if _, ok := plan.Types[PlanType_WRITE]; ok {
			// Get all the nodes that are writable for the given shard ID
			dataNodeShards = tempDataNodeShardIds
			readOnly = false
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
			task.Query,
				task.Type = readPlan.Query, readPlan.Type
		} else if writePlan, ok := plan.Types[PlanType_WRITE]; ok {
			task.Query,
				task.Type = writePlan.Query, writePlan.Type
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
