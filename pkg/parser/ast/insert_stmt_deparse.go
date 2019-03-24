// Auto-generated - DO NOT EDIT

package pg_query

import (
	"fmt"
	"github.com/juju/errors"
	"strings"
)

func (node InsertStmt) Deparse(ctx Context) (*string, error) {
	out := make([]string, 0)
	if node.WithClause != nil {
		if str, err := deparseNode(node.WithClause, Context_None); err != nil {
			return nil, err
		} else {
			out = append(out, *str)
		}
	}

	if node.Relation == nil {
		return nil, errors.New("relation in insert cannot be null!")
	}
	out = append(out, "INSERT INTO")
	if str, err := deparseNode(*node.Relation, Context_None); err != nil {
		return nil, err
	} else {
		out = append(out, *str)
	}

	if node.Cols.Items != nil {
		cols := make([]string, len(node.Cols.Items))
		for i, col := range node.Cols.Items {
			if str, err := deparseNode(col, Context_None); err != nil {
				return nil, err
			} else {
				cols[i] = *str
			}
		}
		out = append(out, fmt.Sprintf("(%s)", strings.Join(cols, ", ")))
	}

	if str, err := node.SelectStmt.Deparse(Context_None); err != nil {
		return nil, err
	} else {
		out = append(out, *str)
	}

	if node.ReturningList.Items != nil && len(node.ReturningList.Items) > 0 {
		out = append(out, "RETURNING")
		fields := make([]string, len(node.ReturningList.Items))
		for i, field := range node.ReturningList.Items {
			if str, err := deparseNode(field, Context_Select); err != nil {
				return nil, err
			} else {
				fields[i] = *str
			}
		}
		out = append(out, strings.Join(fields, ", "))
	}

	result := strings.Join(out, " ")
	return &result, nil
}
