// Auto-generated - DO NOT EDIT

package ast

import (
	"fmt"
	"strings"
)

func (node Alias) Deparse(ctx Context) (string, error) {
	if node.Colnames.Items != nil && len(node.Colnames.Items) > 0 {
		if colnames, err := deparseNodeList(node.Colnames.Items, Context_None); err != nil {
			return "", err
		} else {
			cols := strings.Join(colnames, ", ")
			return fmt.Sprintf(`%s (%s)`, *node.Aliasname, cols), nil
		}
	} else {
		return *node.Aliasname, nil
	}
}
