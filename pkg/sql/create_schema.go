package sql

import (
	"github.com/elliotcourant/noahdb/pkg/ast"
	"github.com/elliotcourant/noahdb/pkg/core"
)

type createSchemaStmtPlanner struct {
	tables []core.Table
	tree   ast.CreateSchemaStmt
}

func newCreateSchemaStatementPlan(tree ast.CreateSchemaStmt) *createSchemaStmtPlanner {
	return &createSchemaStmtPlanner{
		tree: tree,
	}
}

func (stmt *createSchemaStmtPlanner) getNoahQueryPlan(s *session) (InitialPlan, bool, error) {
	return InitialPlan{}, false, nil
}

func (stmt *createSchemaStmtPlanner) GetQueryPlan(s *session) (InitialPlan, bool, error) {
	return InitialPlan{}, false, nil
}
