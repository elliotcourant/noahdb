package planner

import (
	"github.com/elliotcourant/noahdb/pkg/ast"
)

type PlanType string

const (
	READ      PlanType = "READ"
	WRITE     PlanType = "WRITE"
	READWRITE PlanType = "READWRITE"
)

type QueryPlan struct {
	Plans map[PlanType]QueryTask
}

type QueryTask struct {
	Query string
	Type  ast.StmtType
}

type ExecutionPlan struct {
	PlanType PlanType
	Query    string
	Type     ast.StmtType
}
