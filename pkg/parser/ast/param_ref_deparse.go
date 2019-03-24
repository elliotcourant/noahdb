// Auto-generated - DO NOT EDIT

package ast

import (
	"fmt"
)

func (node ParamRef) Deparse(ctx Context) (string, error) {
	result := fmt.Sprintf("$%d", node.Number)
	return &result, nil
}
