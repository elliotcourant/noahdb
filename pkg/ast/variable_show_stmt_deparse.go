// Auto-generated - DO NOT EDIT

package ast

import (
	"fmt"
)

func (node VariableShowStmt) Deparse(ctx Context) (string, error) {
	return fmt.Sprintf("SHOW %s", *node.Name), nil
}
