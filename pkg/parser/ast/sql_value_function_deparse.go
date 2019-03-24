// Auto-generated - DO NOT EDIT

package pg_query

import (
	"fmt"
)

func (node SQLValueFunction) Deparse(ctx Context) (*string, error) {
	switch node.Op {
	case SVFOP_CURRENT_TIMESTAMP:
		result := "CURRENT_TIMESTAMP"
		return &result, nil
	case SVFOP_CURRENT_USER:
		result := "CURRENT_USER"
		return &result, nil
	default:
		panic(fmt.Sprintf("cannot deparse SQLValueFunction %d at this time", node.Op))
	}
	return nil, nil
}
