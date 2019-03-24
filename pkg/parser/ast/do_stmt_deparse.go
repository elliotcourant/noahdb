// Auto-generated - DO NOT EDIT

package ast

import (
	"fmt"
)

func (node DoStmt) Deparse(ctx Context) (string, error) {
	return fmt.Sprintf("DO $$%s$$", node.Args.Items[0].(DefElem).Arg.(String).Str), nil
}
