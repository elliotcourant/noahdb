// Auto-generated - DO NOT EDIT

package pg_query

import (
	"strings"
)

func (node ColumnDef) Deparse(ctx Context) (*string, error) {
	out := []string{*node.Colname}

	if str, err := deparseNode(*node.TypeName, Context_None); err != nil {
		return nil, err
	} else {
		out = append(out, *str)
	}

	if node.RawDefault != nil {
		out = append(out, "USING")
		if str, err := deparseNode(node.RawDefault, Context_None); err != nil {
			return nil, err
		} else {
			out = append(out, *str)
		}
	}

	if node.Constraints.Items != nil && len(node.Constraints.Items) > 0 {
		constraints := make([]string, len(node.Constraints.Items))
		for i, constraint := range node.Constraints.Items {
			if str, err := constraint.Deparse(Context_None); err != nil {
				return nil, err
			} else {
				constraints[i] = *str
			}
		}
		out = append(out, constraints...)
	}
	result := strings.Join(out, " ")
	return &result, nil
}
