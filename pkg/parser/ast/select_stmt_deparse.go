// Auto-generated - DO NOT EDIT

package ast

import (
	"fmt"
	"strings"
)

func (node SelectStmt) Deparse(ctx Context) (string, error) {
	out := make([]string, 0)
	if node.Op == SETOP_UNION {
		if str, err := node.Larg.Deparse(Context_None); err != nil {
			return "", err
		} else {
			out = append(out, str)
		}

		out = append(out, "UNION")
		if node.All {
			out = append(out, "ALL")
		}

		if str, err := node.Rarg.Deparse(Context_None); err != nil {
			return "", err
		} else {
			out = append(out, str)
		}

		return strings.Join(out, " "), nil
	}

	if node.WithClause != nil {
		if str, err := node.WithClause.Deparse(Context_None); err != nil {
			return "", err
		} else {
			out = append(out, str)
		}
	}

	// Get select *distinct* *fields*
	if node.TargetList.Items != nil && len(node.TargetList.Items) > 0 {
		out = append(out, "SELECT")
		if node.DistinctClause.Items != nil && len(node.DistinctClause.Items) > 0 {
			out = append(out, "DISTINCT")
		}
		fields := make([]string, len(node.TargetList.Items))
		for i, field := range node.TargetList.Items {
			if str, err := field.Deparse(Context_Select); err != nil {
				return "", err
			} else {
				fields[i] = str
			}
		}
		out = append(out, strings.Join(fields, ", "))
	}

	if node.FromClause.Items != nil && len(node.FromClause.Items) > 0 {
		out = append(out, "FROM")
		froms := make([]string, len(node.FromClause.Items))
		for i, from := range node.FromClause.Items {
			if str, err := from.Deparse(Context_Select); err != nil {
				return "", err
			} else {
				froms[i] = str
			}
		}
		out = append(out, strings.Join(froms, ", "))
	}

	if node.WhereClause != nil {
		out = append(out, "WHERE")
		if str, err := node.WhereClause.Deparse(Context_None); err != nil {
			return "", err
		} else {
			out = append(out, str)
		}
	}

	if node.ValuesLists != nil && len(node.ValuesLists) > 0 {
		out = append(out, "VALUES")
		allValues := make([]string, len(node.ValuesLists))
		for v, valueList := range node.ValuesLists {
			values := make([]string, len(valueList))
			for i, value := range valueList {
				if str, err := value.Deparse(Context_None); err != nil {
					return "", err
				} else {
					values[i] = str
				}
			}
			allValues[v] = fmt.Sprintf("(%s)", strings.Join(values, ", "))
		}
		out = append(out, strings.Join(allValues, ", "))
	}

	if node.GroupClause.Items != nil && len(node.GroupClause.Items) > 0 {
		out = append(out, "GROUP BY")
		groups := make([]string, len(node.GroupClause.Items))
		for i, group := range node.GroupClause.Items {
			if str, err := group.Deparse(Context_None); err != nil {
				return "", err
			} else {
				groups[i] = str
			}
		}
		out = append(out, strings.Join(groups, ", "))
	}

	if node.HavingClause != nil {
		if str, err := node.HavingClause.Deparse(Context_None); err != nil {
			return "", err
		} else {
			out = append(out, str)
		}
	}

	// Sort clause
	if len(node.SortClause.Items) > 0 {
		out = append(out, "ORDER BY")
		sort := make([]string, len(node.SortClause.Items))
		for i, item := range node.SortClause.Items {
			if str, err := item.Deparse(Context_None); err != nil {
				return "", err
			} else {
				sort[i] = str
			}
		}

		out = append(out, strings.Join(sort, ", "))
	}

	if node.LimitCount != nil {
		out = append(out, "LIMIT")
		if str, err := node.LimitCount.Deparse(Context_None); err != nil {
			return "", err
		} else {
			out = append(out, str)
		}
	}

	if node.LimitOffset != nil {
		out = append(out, "OFFSET")
		if str, err := node.LimitOffset.Deparse(Context_None); err != nil {
			return "", err
		} else {
			out = append(out, str)
		}
	}

	if node.LockingClause.Items != nil && len(node.LockingClause.Items) > 0 {
		for _, lock := range node.LockingClause.Items {
			if str, err := lock.Deparse(Context_None); err != nil {
				return "", err
			} else {
				out = append(out, str)
			}
		}
	}

	return strings.Join(out, " "), nil
}
