// Auto-generated - DO NOT EDIT

package pg_query

import (
	"fmt"
	"strings"
)

func (node Alias) Deparse(ctx Context) (*string, error) {
	if node.Colnames.Items != nil && len(node.Colnames.Items) > 0 {
		if colnames, err := deparseNodeList(node.Colnames.Items, Context_None); err != nil {
			return nil, err
		} else {
			cols := strings.Join(colnames, ", ")
			result := fmt.Sprintf(`%s (%s)`, *node.Aliasname, cols)
			return &result, nil
		}
	} else {
		return node.Aliasname, nil
	}
}
