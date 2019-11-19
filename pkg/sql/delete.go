package sql

import (
	"github.com/elliotcourant/noahdb/pkg/ast"
)

type deleteStmtPlanner struct {
	tree ast.DeleteStmt
}

func newDeleteStatementPlan(tree ast.DeleteStmt) *deleteStmtPlanner {
	return &deleteStmtPlanner{
		tree: tree,
	}
}

func (stmt *deleteStmtPlanner) getNormalQueryPlan(s *session) (InitialPlan, bool, error) {
	return InitialPlan{}, false, nil
}
