package queryutil

import (
	"encoding/json"
	"github.com/elliotcourant/noahdb/pkg/ast"
	"github.com/elliotcourant/timber"
	"testing"
)

var builtInFunctions = BuiltInFunctionMap{
	"pg_catalog.current_database": func(args ...ast.Node) (i interface{}, e error) {
		return "noah", nil
	},
	"pg_catalog.current_schemas": func(args ...ast.Node) (i interface{}, e error) {
		return []string{"noah", "test"}, nil
	},
	"pg_catalog.current_user": func(args ...ast.Node) (i interface{}, e error) {
		return "ecourant", nil
	},
	"pg_catalog.number_of_databases": func(args ...ast.Node) (i interface{}, e error) {
		return 1, nil
	},
}

func Test_ReplaceFunctions(t *testing.T) {
	input := `select current_database(), current_schemas(), CURRENT_USER, number_of_databases();`
	parsed, err := ast.Parse(input)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	stmt := parsed.Statements[0].(ast.RawStmt).Stmt
	stmtJson, _ := json.Marshal(stmt)

	timber.Debugf("INPUT: | %s", string(stmtJson))
	newStmt, err := ReplaceFunctionCalls(stmt, builtInFunctions)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	newStmtJson, _ := json.Marshal(newStmt)
	timber.Debugf("OUTPUT: | %s", string(newStmtJson))
	compiled, err := newStmt.(ast.Node).Deparse(ast.Context_None)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	timber.Infof("IN  | %s", input)
	timber.Infof("OUT | %s", compiled)
}
