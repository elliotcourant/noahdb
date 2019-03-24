// Auto-generated - DO NOT EDIT

package pg_query

import (
	"github.com/juju/errors"
	"strings"
)

func (node UpdateStmt) Deparse(ctx Context) (*string, error) {
	out := make([]string, 0)

	if node.WithClause != nil {
		if str, err := node.WithClause.Deparse(Context_None); err != nil {
			return nil, err
		} else {
			out = append(out, *str)
		}
	}

	out = append(out, "UPDATE")

	if node.Relation == nil {
		return nil, errors.New("relation of update statement cannot be null")
	}

	if str, err := deparseNode(*node.Relation, Context_None); err != nil {
		return nil, err
	} else {
		out = append(out, *str)
	}

	if node.TargetList.Items == nil || len(node.TargetList.Items) == 0 {
		return nil, errors.New("update statement cannot have no sets")
	}

	if node.TargetList.Items != nil && len(node.TargetList.Items) > 0 {
		out = append(out, "SET")
		for _, target := range node.TargetList.Items {
			if str, err := target.Deparse(Context_Update); err != nil {
				return nil, err
			} else {
				out = append(out, *str)
			}
		}
	}

	if node.WhereClause != nil {
		out = append(out, "WHERE")
		if str, err := deparseNode(node.WhereClause, Context_None); err != nil {
			return nil, err
		} else {
			out = append(out, *str)
		}
	}

	if node.ReturningList.Items != nil && len(node.ReturningList.Items) > 0 {
		out = append(out, "RETURNING")
		returning := make([]string, len(node.ReturningList.Items))
		for i, slct := range node.ReturningList.Items {
			if str, err := deparseNode(slct, Context_Select); err != nil {
				return nil, err
			} else {
				returning[i] = *str
			}
		}
		out = append(out, strings.Join(returning, ", "))
	}

	result := strings.Join(out, " ")
	return &result, nil
}
