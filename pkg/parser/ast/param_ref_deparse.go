// Auto-generated - DO NOT EDIT

package ast

import (
	"fmt"
)

func (node ParamRef) Deparse(ctx Context) (string, error) {
	return fmt.Sprintf("$%d", node.Number), nil
}
