package ast

type SyntaxTree struct {
	Statements []Node
	Query      string
}

func (tree *SyntaxTree) Deparse(ctx Context) (string, error) {
	panic("Not Implemented")
}
