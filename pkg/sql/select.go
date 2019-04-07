package sql

import (
	"github.com/elliotcourant/noahdb/pkg/ast"
	"github.com/elliotcourant/noahdb/pkg/util/queryutil"
)

type selectStmtPlanner struct {
	tree ast.SelectStmt
}

func CreateSelectStatementPlan(tree ast.SelectStmt) *selectStmtPlanner {
	return &selectStmtPlanner{
		tree: tree,
	}
}

func (stmt *selectStmtPlanner) getNoahQueryPlan(s *session) (InitialPlan, bool, error) {
	tables := queryutil.GetTables(stmt.tree)
	if len(tables) == 0 {
		return InitialPlan{}, false, nil
	}
	return InitialPlan{}, false, nil
}

func (stmt *selectStmtPlanner) getSimpleQueryPlan(s *session) (InitialPlan, bool, error) {
	tables := queryutil.GetTables(stmt.tree)
	if len(tables) == 0 {
		return InitialPlan{
			Target: PlanTarget_INTERNAL,
			Types: map[PlanType]InitialPlanTask{
				PlanType_READ: {
					Type:  stmt.tree.StatementType(),
					Query: "SELECT 1",
				},
			},
		}, true, nil
	}

	return InitialPlan{}, false, nil
}
