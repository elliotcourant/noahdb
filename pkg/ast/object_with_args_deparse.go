// Auto-generated - DO NOT EDIT

package ast

import (
	"fmt"
	"strings"
)

func (node ObjectWithArgs) Deparse(ctx Context) (string, error) {
	out := make([]string, 0)

	objName := make([]string, len(node.Objname.Items))
	for i, name := range node.Objname.Items {
		if str, err := name.Deparse(Context_FuncCall); err != nil {
			return "", err
		} else {
			objName[i] = str
		}
	}

	args := make([]string, len(node.Objargs.Items))
	for i, arg := range node.Objargs.Items {
		if str, err := arg.Deparse(Context_FuncCall); err != nil {
			return "", err
		} else {
			args[i] = str
		}
	}

	out = append(out, fmt.Sprintf("%s(%s)", strings.Join(objName, "."), strings.Join(args, ", ")))

	return strings.Join(out, " "), nil
}
