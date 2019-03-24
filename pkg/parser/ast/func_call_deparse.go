// Auto-generated - DO NOT EDIT

package pg_query

import (
	"fmt"
	"strings"
)

func (node FuncCall) Deparse(ctx Context) (*string, error) {
	out := make([]string, 0)

	args := make([]string, len(node.Args.Items))
	args, err := deparseNodeList(node.Args.Items, Context_None)
	if err != nil {
		return nil, err
	}

	if node.AggStar {
		args = append(args, "*")
	}

	funcName, err := node.Name()
	if err != nil {
		return nil, err
	}

	distinct := ""
	if node.AggDistinct {
		distinct = "DISTINCT "
	}

	out = append(out, fmt.Sprintf("%s(%s%s)", funcName, distinct, strings.Join(args, ", ")))

	if node.Over != nil {
		if over, err := deparseNode(node.Over, Context_None); err != nil {
			return nil, err
		} else {
			out = append(out, fmt.Sprintf("OVER (%s)", *over))
		}
	}

	result := strings.Join(out, " ")
	return &result, nil
}
