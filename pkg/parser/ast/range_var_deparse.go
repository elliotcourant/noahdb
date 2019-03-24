// Auto-generated - DO NOT EDIT

package pg_query

import (
	"fmt"
	"strings"
)

func (node RangeVar) Deparse(ctx Context) (*string, error) {
	out := make([]string, 0)
	if !node.Inh {
		out = append(out, "ONLY")
	}

	if node.Schemaname != nil && len(*node.Schemaname) > 0 {
		out = append(out, fmt.Sprintf(`"%s"."%s"`, *node.Schemaname, *node.Relname))
	} else {
		out = append(out, fmt.Sprintf(`"%s"`, *node.Relname))
	}

	if node.Alias != nil {
		if str, err := deparseNode(*node.Alias, Context_None); err != nil {
			return nil, err
		} else {
			out = append(out, *str)
		}
	}

	result := strings.Join(out, " ")
	return &result, nil
}
