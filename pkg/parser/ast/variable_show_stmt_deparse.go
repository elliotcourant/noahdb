// Auto-generated - DO NOT EDIT

package pg_query

import (
	"strings"
)

func (node VariableShowStmt) Deparse(ctx Context) (*string, error) {
	out := []string{"SHOW"}
	out = append(out, *node.Name)
	result := strings.Join(out, " ")
	return &result, nil
}
