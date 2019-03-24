// Auto-generated - DO NOT EDIT

package pg_query

import (
	"fmt"
	"github.com/juju/errors"
	"strings"
)

func (node JoinExpr) Deparse(ctx Context) (*string, error) {
	out := make([]string, 0)

	if node.Larg == nil {
		return nil, errors.New("larg of join cannot be null")
	}

	if str, err := deparseNode(node.Larg, Context_None); err != nil {
		return nil, err
	} else {
		out = append(out, *str)
	}

	switch node.Jointype {
	case JOIN_INNER:
		if node.IsNatural {
			out = append(out, "NATURAL")
		} else if node.Quals == nil && (node.UsingClause.Items == nil || len(node.UsingClause.Items) == 0) {
			out = append(out, "CROSS")
		}
	case JOIN_LEFT:
		out = append(out, "LEFT")
	case JOIN_FULL:
		out = append(out, "FULL")
	case JOIN_RIGHT:
		out = append(out, "RIGHT")
	default:
		return nil, errors.Errorf("cannot handle join type (%d)", node.Jointype)
	}
	out = append(out, "JOIN")

	if node.Rarg == nil {
		return nil, errors.New("rarg of join cannot be null")
	}

	if str, err := deparseNode(node.Rarg, Context_None); err != nil {
		return nil, err
	} else {
		out = append(out, *str)
	}

	if node.Quals != nil {
		out = append(out, "ON")
		if str, err := deparseNode(node.Quals, Context_None); err != nil {
			return nil, err
		} else {
			out = append(out, *str)
		}
	}

	if node.UsingClause.Items != nil && len(node.UsingClause.Items) > 0 {
		clauses := make([]string, len(node.UsingClause.Items))
		for i, field := range node.UsingClause.Items {
			if str, err := deparseNode(field, Context_Select); err != nil {
				return nil, err
			} else {
				clauses[i] = *str
			}
		}
		out = append(out, fmt.Sprintf("USING (%s)", strings.Join(clauses, ", ")))
	}

	result := strings.Join(out, " ")
	return &result, nil
}
