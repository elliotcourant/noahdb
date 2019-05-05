package sql

import (
	"github.com/elliotcourant/noahdb/pkg/ast"
	"github.com/elliotcourant/noahdb/pkg/core"
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
	Tasks []ExpandedPlanTask
}

type ExpandedPlanTask struct {
	Query    string
	ReadOnly bool
	DataNode core.DataNode
	Shard    core.Shard
	Type     ast.StmtType
}

func (s *session) expandQueryPlan(plan InitialPlan) (ExpandedPlan, error) {
	startTimestamp := time.Now()
	defer func() {
		golog.Verbosef("expanding of plan took %s", time.Since(startTimestamp))
	}()

	if plan.Target == PlanTarget_INTERNAL {
		// Internal query plans can go directly to the SQLite database.
		golog.Verbosef("plan targets internal SQLite database")
		panic("internal queries are not yet supported")
	}

	nodes := make([]core.DataNode, 0)
	switch plan.ShardID {
	case 0: // If this query does not target a specific shard.
		if _, ok := plan.Types[PlanType_READ]; ok {
			// Get a single node to execute the read query.
		}
	default:
		if _, ok := plan.Types[PlanType_READ]; ok {
			// Get any readable node for the given shard ID
		}
	}

	return ExpandedPlan{
		Tasks: make([]ExpandedPlanTask, 0),
	}, nil
}
