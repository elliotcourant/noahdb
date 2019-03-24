// Auto-generated - DO NOT EDIT

package ast

import (
	"fmt"
	"strings"
)

func (node SubLink) Deparse(ctx Context) (string, error) {
	switch node.SubLinkType {
	case EXPR_SUBLINK:
		if subSelect, err := node.Subselect.Deparse(Context_None); err != nil {
			return "", err
		} else {
			return fmt.Sprintf("(%s)", subSelect), err
		}
	case ANY_SUBLINK:
		out := []string{"", "IN", ""}
		if columnRef, err := node.Testexpr.Deparse(Context_None); err != nil {
			return "", err
		} else {
			out[0] = columnRef
		}

		if subSelect, err := node.Subselect.Deparse(Context_None); err != nil {
			return "", err
		} else {
			out[2] = fmt.Sprintf("(%s)", subSelect)
		}

		return strings.Join(out, " "), nil
	default:
		panic(fmt.Sprintf("cannot handle sub link type [%s]", node.SubLinkType.String()))
	}
	return "", nil
}
