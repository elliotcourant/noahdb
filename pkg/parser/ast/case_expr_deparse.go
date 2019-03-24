// Auto-generated - DO NOT EDIT

package pg_query

import (
	"github.com/juju/errors"
	"strings"
)

func (node CaseExpr) Deparse(ctx Context) (*string, error) {
	out := []string{"CASE"}

	if node.Arg != nil {
		if str, err := deparseNode(node.Arg, Context_None); err != nil {
			return nil, err
		} else {
			out = append(out, *str)
		}
	}

	if node.Args.Items == nil || len(node.Args.Items) == 0 {
		return nil, errors.New("case expression cannot have 0 arguments")
	}

	if args, err := deparseNodeList(node.Args.Items, Context_None); err != nil {
		return nil, err
	} else {
		out = append(out, args...)
	}

	if node.Defresult != nil {
		out = append(out, "ELSE")
		if str, err := deparseNode(node.Defresult, Context_None); err != nil {
			return nil, err
		} else {
			out = append(out, *str)
		}
	}

	out = append(out, "END")
	result := strings.Join(out, " ")
	return &result, nil
}
