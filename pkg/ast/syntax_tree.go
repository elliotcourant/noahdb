package ast

import (
	"encoding/json"
)

type SyntaxTree struct {
	Statements []Node
	Query      string
}

func (tree *SyntaxTree) Deparse(ctx Context) (string, error) {
	panic("Not Implemented")
}

func (input SyntaxTree) MarshalJSON() ([]byte, error) {
	type ParsetreeListAlias SyntaxTree
	return json.Marshal(input.Statements)
}

func (output *SyntaxTree) UnmarshalJSON(input []byte) (err error) {
	var list []json.RawMessage
	err = json.Unmarshal([]byte(input), &list)
	if err != nil {
		return
	}

	for _, nodeJson := range list {
		var node Node
		node, err = UnmarshalNodeJSON(nodeJson)
		if err != nil {
			return
		}
		output.Statements = append(output.Statements, node)
	}

	return
}
