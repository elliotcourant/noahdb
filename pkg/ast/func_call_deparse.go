// Auto-generated - DO NOT EDIT

package ast

import (
	"fmt"
	"strings"
)

func (node FuncCall) Deparse(ctx Context) (string, error) {
	out := make([]string, 0)

	args := make([]string, len(node.Args.Items))
	args, err := deparseNodeList(node.Args.Items, Context_None)
	if err != nil {
		return "", err
	}

	if node.AggStar {
		args = append(args, "*")
	}

	funcName, err := node.Name()
	if err != nil {
		return "", err
	}

	distinct := ""
	if node.AggDistinct {
		distinct = "DISTINCT "
	}

	out = append(out, fmt.Sprintf("%s(%s%s)", funcName, distinct, strings.Join(args, ", ")))

	if node.Over != nil {
		if over, err := node.Over.Deparse(Context_None); err != nil {
			return "", err
		} else {
			out = append(out, fmt.Sprintf("OVER (%s)", over))
		}
	}

	return strings.Join(out, " "), nil
}
