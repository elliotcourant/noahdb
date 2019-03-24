// Auto-generated - DO NOT EDIT

package ast

import (
	"github.com/juju/errors"
	"strings"
)

func (node CaseExpr) Deparse(ctx Context) (string, error) {
	out := []string{"CASE"}

	if node.Arg != nil {
		if str, err := node.Arg.Deparse(Context_None); err != nil {
			return "", err
		} else {
			out = append(out, str)
		}
	}

	if node.Args.Items == nil || len(node.Args.Items) == 0 {
		return "", errors.New("case expression cannot have 0 arguments")
	}

	if args, err := deparseNodeList(node.Args.Items, Context_None); err != nil {
		return "", err
	} else {
		out = append(out, args...)
	}

	if node.Defresult != nil {
		out = append(out, "ELSE")
		if str, err := node.Defresult.Deparse(Context_None); err != nil {
			return "", err
		} else {
			out = append(out, str)
		}
	}

	out = append(out, "END")
	return strings.Join(out, " "), nil
}
