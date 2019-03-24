// Auto-generated - DO NOT EDIT

package pg_query

import (
	"strings"
)

func (node VariableSetStmt) Deparse(ctx Context) (*string, error) {
	out := []string{"SET"}
	if node.IsLocal {
		out = append(out, "LOCAL")
	}
	out = append(out, *node.Name)
	out = append(out, "TO")
	if args, err := deparseNodeList(node.Args.Items, Context_None); err != nil {
		return nil, err
	} else {
		out = append(out, args...)
	}
	result := strings.Join(out, " ")
	return &result, nil
}
