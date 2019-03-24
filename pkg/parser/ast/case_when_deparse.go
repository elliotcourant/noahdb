// Auto-generated - DO NOT EDIT

package ast

import (
	"strings"
)

func (node CaseWhen) Deparse(ctx Context) (string, error) {
	// The 1st blank string will be replaced by node.Expr
	// The 2nd blank string will be replaced by node.Result
	out := []string{"WHEN", "", "THEN", ""}

	if str, err := node.Expr.Deparse(Context_None); err != nil {
		return "", err
	} else {
		out[1] = str
	}

	if str, err := node.Result.Deparse(Context_None); err != nil {
		return "", err
	} else {
		out[3] = str
	}

	return strings.Join(out, " "), nil
}
