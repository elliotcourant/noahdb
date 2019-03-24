// Auto-generated - DO NOT EDIT

package pg_query

import (
	"strings"
)

func (node CaseWhen) Deparse(ctx Context) (*string, error) {
	// The 1st blank string will be replaced by node.Expr
	// The 2nd blank string will be replaced by node.Result
	out := []string{"WHEN", "", "THEN", ""}

	if str, err := deparseNode(node.Expr, Context_None); err != nil {
		return nil, err
	} else {
		out[1] = *str
	}

	if str, err := deparseNode(node.Result, Context_None); err != nil {
		return nil, err
	} else {
		out[3] = *str
	}

	result := strings.Join(out, " ")
	return &result, nil
}
