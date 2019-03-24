// Auto-generated - DO NOT EDIT

package ast

import (
	"fmt"
	"strings"
)

func (node String) Deparse(ctx Context) (string, error) {
	switch ctx {
	case Context_AConst:
		return fmt.Sprintf("'%s'", strings.Replace(node.Str, "'", "''", -1)), nil
	case Context_FuncCall, Context_TypeName, Context_Operator:
		return node.Str, nil
	default:
		return fmt.Sprintf(`"%s"`, strings.Replace(node.Str, `"`, `""`, -1)), nil
	}
}
