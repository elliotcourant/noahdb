// Auto-generated - DO NOT EDIT

package pg_query

import (
	"fmt"
	"strings"
)

func (node ObjectWithArgs) Deparse(ctx Context) (*string, error) {
	out := make([]string, 0)

	objName := make([]string, len(node.Objname.Items))
	for i, name := range node.Objname.Items {
		if str, err := name.Deparse(Context_FuncCall); err != nil {
			return nil, err
		} else {
			objName[i] = *str
		}
	}

	args := make([]string, len(node.Objargs.Items))
	for i, arg := range node.Objargs.Items {
		if str, err := arg.Deparse(Context_FuncCall); err != nil {
			return nil, err
		} else {
			args[i] = *str
		}
	}

	out = append(out, fmt.Sprintf("%s(%s)", strings.Join(objName, "."), strings.Join(args, ", ")))

	result := strings.Join(out, " ")
	return &result, nil
}
