// Auto-generated - DO NOT EDIT

package ast

import (
	"github.com/juju/errors"
	"strings"
)

func (node ResTarget) Deparse(ctx Context) (*string, error) {
	switch ctx {
	case Context_None:
		return node.Name, nil
	case Context_Select:
		out := make([]string, 0)
		if str, err := node.Val.Deparse(Context_Select); err != nil {
			return nil, err
		} else {
			out = append(out, *str)
		}

		if node.Name != nil && len(*node.Name) > 0 {
			out = append(out, "AS")
			out = append(out, *node.Name)
		}
		result := strings.Join(out, " ")
		return &result, nil
	case Context_Update:
		out := make([]string, 0)
		if node.Name == nil || len(*node.Name) == 0 {
			return nil, errors.New("cannot have blank name for res target in update")
		}
		out = append(out, *node.Name)

		if node.Val == nil {
			return nil, errors.New("cannot have null value for res target in update")
		}

		if str, err := deparseNode(node.Val, Context_None); err != nil {
			return nil, err
		} else {
			out = append(out, *str)
		}

		result := strings.Join(out, " = ")
		return &result, nil
	default:
		return nil, errors.Errorf("context type %s is not currently implemented", ctx.String())
	}
}
