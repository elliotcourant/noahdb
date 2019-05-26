package sql

import (
	"fmt"
	"github.com/elliotcourant/noahdb/pkg/ast"
)

type createStmtPlanner struct {
	tree ast.CreateStmt
}

func NewCreateStatementPlan(tree ast.CreateStmt) *createStmtPlanner {
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
	return InitialPlan{}, false, nil
}
