// Auto-generated - DO NOT EDIT

package pg_query

import (
	"fmt"
)

func (node DefElem) Deparse(ctx Context) (*string, error) {
	if arg, err := node.Arg.Deparse(Context_AConst); err != nil {
		return nil, err
	} else {
		result := fmt.Sprintf("%s %s", *node.Defname, *arg)
		return &result, nil
	}
}
