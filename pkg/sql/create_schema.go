package sql

import (
	"github.com/elliotcourant/noahdb/pkg/ast"
	"github.com/elliotcourant/noahdb/pkg/core"
)

type createSchemaStmtPlanner struct {
	tables []core.Table
	tree   ast.SelectStmt
}

func CreateCreateSchemaStatementPlan(tree ast.CreateSchemaStmt) *selectStmtPlanner {
	return &selectStmtPlanner{
		tree: tree,
	}
}

func (stmt *createSchemaStmtPlanner) getNoahQueryPlan(s *session) (InitialPlan, bool, error) {
	return InitialPlan{}, false, nil
}

func (stmt *createSchemaStmtPlanner) getSimpleQueryPlan(s *session) (InitialPlan, bool, error) {
	return InitialPlan{}, false, nil
}
