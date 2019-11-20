package sql

import (
	"github.com/elliotcourant/noahdb/pkg/ast"
	"github.com/stretchr/testify/assert"
	"testing"
)

func GetStatementPlanner(t *testing.T, query string) QueryPlanner {
	tree, err := ast.Parse(query)
	if !assert.NoError(t, err) {
		panic(err)
	}

	planner, err := getStatementHandler(tree.Statements[0].(ast.RawStmt).Stmt.(ast.Stmt))
	if !assert.NoError(t, err) {
		panic(err)
	}

	return planner.(QueryPlanner)
}
