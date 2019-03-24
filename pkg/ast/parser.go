package ast

import (
	"encoding/json"
	"github.com/elliotcourant/noahdb/pkg/ast/pgparser"
	"runtime/debug"
)

// ParseToJSON - Parses the given SQL statement into an AST (JSON format)
func parseToJSON(input string) (result string, err error) {
	return pgparser.ParseToJSON(input)
}

// Parse the given SQL statement into an AST (native Go structs)
func Parse(input string) (tree SyntaxTree, err error) {
	defer func() {
		if r := recover(); r != nil {
			debug.PrintStack()
			err = r.(error)
		}
	}()
	jsonTree, err := parseToJSON(input)
	if err != nil {
		return
	}

	// JSON unmarshalling can panic in edge cases we don't support yet. This is
	// still a *bug that needs to be fixed*, but this way the caller can expect an
	// error to be returned always, instead of a panic

	err = json.Unmarshal([]byte(jsonTree), &tree)
	tree.Query = input
	return
}
