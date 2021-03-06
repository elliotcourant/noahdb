// Auto-generated - DO NOT EDIT

package ast

import (
	"fmt"
	"github.com/juju/errors"
	"strings"
)

func (node InsertStmt) Deparse(ctx Context) (string, error) {
	out := make([]string, 0)
	if node.WithClause != nil {
		if str, err := node.WithClause.Deparse(Context_None); err != nil {
			return "", err
		} else {
			out = append(out, str)
		}
	}

	if node.Relation == nil {
		return "", errors.New("relation in insert cannot be null!")
	}
	out = append(out, "INSERT INTO")
	if str, err := (*node.Relation).Deparse(Context_None); err != nil {
		return "", err
	} else {
		out = append(out, str)
	}

	if node.Cols.Items != nil {
		cols := make([]string, len(node.Cols.Items))
		for i, col := range node.Cols.Items {
			if str, err := col.Deparse(Context_None); err != nil {
				return "", err
			} else {
				cols[i] = str
			}
		}
		out = append(out, fmt.Sprintf("(%s)", strings.Join(cols, ", ")))
	}

	if str, err := node.SelectStmt.Deparse(Context_None); err != nil {
		return "", err
	} else {
		out = append(out, str)
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
