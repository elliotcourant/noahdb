// Auto-generated - DO NOT EDIT

package ast

import (
	"github.com/juju/errors"
	"strings"
)

func (node ResTarget) Deparse(ctx Context) (string, error) {
	switch ctx {
	case Context_None:
		return *node.Name, nil
	case Context_Select:
		out := make([]string, 0)
		if str, err := node.Val.Deparse(Context_Select); err != nil {
			return "", err
		} else {
			out = append(out, str)
		}

		if node.Name != nil && len(*node.Name) > 0 {
			out = append(out, "AS")
			out = append(out, *node.Name)
		}
		return strings.Join(out, " "), nil
	case Context_Update:
		out := make([]string, 0)
		if node.Name == nil || len(*node.Name) == 0 {
			return "", errors.New("cannot have blank name for res target in update")
		}
		out = append(out, *node.Name)

		if node.Val == nil {
			return "", errors.New("cannot have null value for res target in update")
		}

		if str, err := node.Val.Deparse(Context_None); err != nil {
			return "", err
		} else {
			out = append(out, str)
		}

		return strings.Join(out, " = "), nil
	default:
		return "", errors.Errorf("context type %s is not currently implemented", ctx.String())
	}
}
