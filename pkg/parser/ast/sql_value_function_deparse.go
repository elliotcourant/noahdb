// Auto-generated - DO NOT EDIT

package ast

import (
	"fmt"
)

func (node SQLValueFunction) Deparse(ctx Context) (string, error) {
	switch node.Op {
	case SVFOP_CURRENT_TIMESTAMP:
		return "CURRENT_TIMESTAMP", nil
	case SVFOP_CURRENT_USER:
		return "CURRENT_USER", nil
	default:
		panic(fmt.Sprintf("cannot deparse SQLValueFunction %s at this time", node.Op.String()))
	}
	return "", nil
}
