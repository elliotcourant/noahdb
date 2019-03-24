// Auto-generated - DO NOT EDIT

package ast

import (
	"strings"
)

func (node DeleteStmt) Deparse(ctx Context) (string, error) {
	out := make([]string, 0)
	if node.WithClause != nil {
		if str, err := node.WithClause.Deparse(Context_None); err != nil {
			return "", err
		} else {
			out = append(out, str)
		}
	}

	out = append(out, "DELETE FROM")

	if table, err := node.Relation.Deparse(Context_None); err != nil {
		return "", err
	} else {
		out = append(out, table)
	}

	if len(node.UsingClause.Items) > 0 {
		out = append(out, "USING")
		using := make([]string, len(node.UsingClause.Items))
		for i, usingItem := range node.UsingClause.Items {
			if str, err := usingItem.Deparse(Context_None); err != nil {
				return "", err
			} else {
				using[i] = str
			}
		}
		out = append(out, strings.Join(using, ", "))
	}

	if node.WhereClause != nil {
		out = append(out, "WHERE")
		if str, err := node.WhereClause.Deparse(Context_None); err != nil {
			return "", err
		} else {
			out = append(out, str)
		}
	}

	if node.ReturningList.Items != nil && len(node.ReturningList.Items) > 0 {
		out = append(out, "RETURNING")
		fields := make([]string, len(node.ReturningList.Items))
		for i, field := range node.ReturningList.Items {
			if str, err := field.Deparse(Context_Select); err != nil {
				return "", err
			} else {
				fields[i] = str
			}
		}
		out = append(out, strings.Join(fields, ", "))
	}

	return strings.Join(out, " "), nil
}
