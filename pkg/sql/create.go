package sql

import (
	"fmt"
	"github.com/elliotcourant/noahdb/pkg/ast"
)

type createStmtPlanner struct {
	tree ast.CreateStmt
}

func NewCreateStatementPlan(tree ast.CreateStmt) *createStmtPlanner {
	// ast.VariableShowStmt{}
	return &createStmtPlanner{
		tree: tree,
	}
}

func (stmt *createStmtPlanner) getSimpleQueryPlan(s *session) (InitialPlan, bool, error) {
	// schemaName := *stmt.tree.Relation.Schemaname
	tableName := *stmt.tree.Relation.Relname
	tables, err := s.Colony().Tables().GetTables(tableName)
	if err != nil {
		return InitialPlan{}, false, fmt.Errorf("could not verify table doesn't exit: %v", err)
	}
	if len(tables) > 0 && !stmt.tree.IfNotExists {
		return InitialPlan{}, false, fmt.Errorf("table with name [%s] already exists", tableName)
	}

	compiledQuery, err := stmt.tree.Deparse(ast.Context_None)
	if err != nil {
		return InitialPlan{}, false, fmt.Errorf("could not recompile query: %v", err)
	}

	return InitialPlan{
		Target:  PlanTarget_STANDARD,
		ShardID: 0,
		Types: map[PlanType]InitialPlanTask{
			PlanType_WRITE: {
				Query: compiledQuery,
				Type:  stmt.tree.StatementType(),
			},
		},
	}, true, nil
}
