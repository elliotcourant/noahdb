// Auto-generated - DO NOT EDIT

package ast

import (
	"fmt"
	"strings"
)

func (node UpdateStmt) Deparse(ctx Context) (string, error) {
	out := make([]string, 0)

	if node.WithClause != nil {
		if str, err := node.WithClause.Deparse(Context_None); err != nil {
			return "", err
		} else {
			out = append(out, str)
		}
	}

	out = append(out, "UPDATE")

	if node.Relation == nil {
		return "", fmt.Errorf("relation of update statement cannot be null")
	}

	if str, err := (*node.Relation).Deparse(Context_None); err != nil {
		return "", err
	} else {
		out = append(out, str)
	}

	if node.TargetList.Items == nil || len(node.TargetList.Items) == 0 {
		return "", fmt.Errorf("update statement cannot have no sets")
	}

	if node.TargetList.Items != nil && len(node.TargetList.Items) > 0 {
		out = append(out, "SET")
		sets := make([]string, len(node.TargetList.Items))
		for i, target := range node.TargetList.Items {
			if str, err := target.Deparse(Context_Update); err != nil {
				return "", err
			} else {
				sets[i] = str
			}
		}
		out = append(out, strings.Join(sets, ", "))
	}

	if len(node.FromClause.Items) > 0 {
		out = append(out, "FROM")
		if strs, err := node.FromClause.DeparseList(Context_Update); err != nil {
			return "", err
		} else {
			out = append(out, strings.Join(strs, ", "))
		}
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
		returning := make([]string, len(node.ReturningList.Items))
		for i, slct := range node.ReturningList.Items {
			if str, err := slct.Deparse(Context_Select); err != nil {
				return "", err
			} else {
				returning[i] = str
			}
		}
		out = append(out, strings.Join(returning, ", "))
	}

	return strings.Join(out, " "), nil
}
