// Auto-generated - DO NOT EDIT

package pg_query

import (
	"strings"
)

func (node DropdbStmt) Deparse(ctx Context) (*string, error) {
	out := []string{"DROP DATABASE"}
	if node.MissingOk {
		out = append(out, "IF EXISTS")
	}
	out = append(out, *node.Dbname)
	result := strings.Join(out, " ")
	return &result, nil
}
