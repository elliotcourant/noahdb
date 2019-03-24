// Auto-generated - DO NOT EDIT

package pg_query

import (
	"fmt"
	"strings"
)

func (node String) Deparse(ctx Context) (*string, error) {
	switch ctx {
	case Context_AConst:
		result := fmt.Sprintf("'%s'", strings.Replace(node.Str, "'", "''", -1))
		return &result, nil
	case Context_FuncCall, Context_TypeName, Context_Operator:
		return &node.Str, nil
	default:
		result := fmt.Sprintf(`"%s"`, strings.Replace(node.Str, `"`, `""`, -1))
		return &result, nil
	}
}
