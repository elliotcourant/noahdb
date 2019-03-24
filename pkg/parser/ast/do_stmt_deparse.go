// Auto-generated - DO NOT EDIT

package ast

import (
	"fmt"
)

func (node DoStmt) Deparse(ctx Context) (*string, error) {
	result := fmt.Sprintf("DO $$%s$$", node.Args.Items[0].(DefElem).Arg.(String).Str)
	return &result, nil
}
