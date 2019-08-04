// Auto-generated - DO NOT EDIT

package ast

import (
	"strings"
)

func (node RangeSubselect) Deparse(ctx Context) (string, error) {
	out := make([]string, 0)
	out = append(out, "(")

	if from, err := node.Subquery.Deparse(Context_Select); err != nil {
		return "", err
	} else {
		out = append(out, from)
	}

	out = append(out, ")")
	if node.Alias != nil {
		out = append(out, "AS")
		out = append(out, *node.Alias.Aliasname)
		if cols, err := node.Alias.Colnames.DeparseList(Context_Select); err != nil {
			return "", err
		} else {
			out = append(out, "(", strings.Join(cols, ", "), ")")
		}
	}

	return strings.Join(out, " "), nil
}
