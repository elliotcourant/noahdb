package sql

import (
	"fmt"
	"github.com/ahmetb/go-linq/v3"
	"github.com/elliotcourant/noahdb/pkg/ast"
	"github.com/elliotcourant/noahdb/pkg/core"
	"github.com/elliotcourant/noahdb/pkg/util/queryutil"
	"strings"
)

type deleteStmtPlanner struct {
	tables []core.Table
	tree   ast.DeleteStmt
}

func newDeleteStatementPlan(tree ast.DeleteStmt) *deleteStmtPlanner {
	return &deleteStmtPlanner{
		tree: tree,
	}
}

func (stmt *deleteStmtPlanner) GetQueryPlan(s *session) (InitialPlan, bool, error) {
	tableNames := queryutil.GetTables(stmt.tree)
	if len(tableNames) == 0 {
		return InitialPlan{}, false, nil
	}

	linq.From(tableNames).Distinct().ToSlice(&tableNames)

	tables, err := s.Colony().Tables().GetTables(tableNames...)
	if err != nil {
		return InitialPlan{}, false, err
	}

	if len(tables) != len(tableNames) {
		// This means that there is a table missing.
		missingTables := make([]string, 0)
		linq.From(tableNames).
			ExceptBy(linq.From(tables), func(i interface{}) interface{} {
				if table, ok := i.(core.Table); ok {
					return table.TableName
				}
				return nil
			}).ToSlice(&missingTables)
		s.log.Debugf("could not resolve tables: %s", strings.Join(missingTables, ", "))
		return InitialPlan{}, false, fmt.Errorf("could not resolve tables with names: %s", strings.Join(missingTables, ", "))
	}

	stmt.tables = tables

	return InitialPlan{}, false, nil
}
