// Auto-generated - DO NOT EDIT

package ast

import (
	"strings"
)

func (node DropdbStmt) Deparse(ctx Context) (string, error) {
	out := []string{"DROP DATABASE"}
	if node.MissingOk {
		out = append(out, "IF EXISTS")
	}
	out = append(out, *node.Dbname)
	return strings.Join(out, " "), nil
}
