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
	Types  map[PlanType]InitialPlanTask
	Target PlanTarget
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

	return ExpandedPlan{
		Tasks: make([]ExpandedPlanTask, 0),
	}, nil
}
