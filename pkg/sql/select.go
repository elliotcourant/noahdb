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
	// If there are no tables then we can simply recompile the query and send it to SQLite,
	// this will make queries like CURRENT_TIMESTAMP or 1 very fast
	if len(tables) == 0 {
		query, err := stmt.tree.Deparse(ast.Context_None)
		if err != nil {
			return InitialPlan{}, true, err
		}
		return InitialPlan{
			Target: PlanTarget_INTERNAL,
			Types: map[PlanType]InitialPlanTask{
				PlanType_READ: {
					Type:  stmt.tree.StatementType(),
					Query: query,
				},
			},
		}, true, nil
	}

	return InitialPlan{}, false, nil
}
