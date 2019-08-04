// Auto-generated - DO NOT EDIT

package ast

import (
	"fmt"
	"strings"
)

func (node CoalesceExpr) Deparse(ctx Context) (string, error) {
	out := make([]string, 0)

	args := make([]string, len(node.Args.Items))
	args, err := deparseNodeList(node.Args.Items, Context_None)
	if err != nil {
		return "", err
	}

	out = append(out, fmt.Sprintf("COALESCE(%s)", strings.Join(args, ", ")))

	return strings.Join(out, " "), nil
}
