package parser

import (
	"github.com/elliotcourant/noahdb/pkg/parser/ast"
)

func Parse(input string) (tree ast.SyntaxTree, err error) {
	return ast.ParseAST(input)
}

func Deparse(node ast.Node) (string, error) {
	return node.Deparse(ast.Context_None)
}
