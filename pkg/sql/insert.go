package sql

import (
	"fmt"
	"github.com/elliotcourant/noahdb/pkg/ast"
	"github.com/elliotcourant/noahdb/pkg/core"
)

type insertStmtPlanner struct {
	tree ast.InsertStmt
}

func newInsertStatementPlan(tree ast.InsertStmt) *insertStmtPlanner {
	return &insertStmtPlanner{
		tree: tree,
	}
}

func (stmt *insertStmtPlanner) getSimpleQueryPlan(s *session) (InitialPlan, bool, error) {
	tableName := *stmt.tree.Relation.Relname
	tables, err := s.Colony().Tables().GetTables(tableName)
	if err != nil {
		return InitialPlan{}, false, err
	}

	// Handle any number of returned tables.
	switch len(tables) {
	case 0:
		return InitialPlan{}, false, fmt.Errorf("could not resolve table [%s]", tableName)
	case 1: // The desired number of tables returned
	default:
		return InitialPlan{}, false, fmt.Errorf("found multiple tables with name [%s]", tableName)
	}

	table := tables[0] // We only want to work with one table.

	switch table.TableType {
	case core.TableType_Noah:
		panic("not handling this yet")
	case core.TableType_Global:

	case core.TableType_Tenant:

	case core.TableType_Sharded:

	}
	return InitialPlan{}, false, nil
}
