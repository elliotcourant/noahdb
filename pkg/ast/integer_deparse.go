// Auto-generated - DO NOT EDIT

package ast

import (
	"strconv"
)

func (node Integer) Deparse(ctx Context) (string, error) {
	return strconv.FormatInt(node.Ival, 10), nil
}
