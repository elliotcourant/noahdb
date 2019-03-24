// Auto-generated - DO NOT EDIT

package pg_query

import (
	"fmt"
	"github.com/juju/errors"
	"reflect"
	"strings"
)

func (node A_Expr) Deparse(ctx Context) (*string, error) {
	switch node.Kind {
	case AEXPR_OP:
		return node.deparseAexpr(ctx)
	case AEXPR_OP_ANY:
		return node.deparseAexprAny(ctx)
	case AEXPR_IN:
		return node.deparseAexprIn(ctx)
	case AEXPR_BETWEEN, AEXPR_NOT_BETWEEN, AEXPR_BETWEEN_SYM, AEXPR_NOT_BETWEEN_SYM:
		return node.deparseAexprBetween(ctx)
	case AEXPR_NULLIF:
		return node.deparseAExprNullIf(ctx)
	default:
		return nil, errors.Errorf("could not parse AExpr of kind: %d, not implemented", node.Kind)
	}
}

func (node A_Expr) deparseAexpr(ctx Context) (*string, error) {
	out := make([]string, 0)
	if node.Lexpr != nil {
		switch n := node.Lexpr.(type) {
		case List:
			if n.Items == nil || len(n.Items) == 0 {
				return nil, errors.New("lexpr list cannot be empty")
			}
			if str, err := n.Items[0].Deparse(ctx); err != nil {
				return nil, err
			} else {
				out = append(out, *str)
			}
		default:
			if str, err := n.Deparse(ctx); err != nil {
				return nil, err
			} else {
				out = append(out, *str)
			}
		}
	}

	if node.Rexpr != nil {
		if str, err := node.Rexpr.Deparse(ctx); err != nil {
			return nil, err
		} else {
			out = append(out, *str)
		}
	}

	if node.Name.Items == nil || len(node.Name.Items) == 0 {
		return nil, errors.New("error, expression name cannot be null")
	}

	if name, err := deparseNode(node.Name.Items[0], Context_Operator); err != nil {
		return nil, err
	} else {
		result := strings.Join(out, fmt.Sprintf(" %s ", *name))
		if ctx != Context_None {
			result = fmt.Sprintf("(%s)", result)
		}
		return &result, nil
	}
}

func (node A_Expr) deparseAexprIn(ctx Context) (*string, error) {
	out := make([]string, 0)

	if node.Rexpr == nil {
		return nil, errors.New("rexpr of IN expression cannot be null")
	}

	// TODO (@elliotcourant) convert to handle list
	if strs, err := deparseNodeList(node.Rexpr.(List).Items, Context_None); err != nil {
		return nil, err
	} else {
		out = append(out, strs...)
	}

	if node.Name.Items == nil || len(node.Name.Items) == 0 {
		return nil, errors.New("names of IN expression cannot be empty")
	}

	if strs, err := deparseNodeList(node.Name.Items, Context_Operator); err != nil {
		return nil, err
	} else {
		operator := ""
		if reflect.DeepEqual(strs, []string{"="}) {
			operator = "IN"
		} else {
			operator = "NOT IN"
		}

		if node.Lexpr == nil {
			return nil, errors.New("lexpr of IN expression cannot be null")
		}

		if str, err := node.Lexpr.Deparse(Context_None); err != nil {
			return nil, err
		} else {
			result := fmt.Sprintf("%s %s (%s)", *str, operator, strings.Join(out, ", "))
			return &result, nil
		}
	}
}

func (node A_Expr) deparseAexprAny(ctx Context) (*string, error) {
	out := make([]string, 0)
	if str, err := node.Lexpr.Deparse(Context_None); err != nil {
		return nil, err
	} else {
		out = append(out, *str)
	}

	if str, err := node.Rexpr.Deparse(Context_None); err != nil {
		return nil, err
	} else {
		out = append(out, fmt.Sprintf("ANY(%s)", *str))
	}

	// TODO (elliotcourant) this seems a bit weird that we are just taking the first item for this. Revist this in the future?
	if str, err := deparseNode(node.Name.Items[0], Context_Operator); err != nil {
		return nil, err
	} else {
		result := strings.Join(out, *str)
		return &result, nil
	}
}

func (node A_Expr) deparseAexprBetween(ctx Context) (*string, error) {
	between := ""
	switch node.Kind {
	case AEXPR_BETWEEN:
		between = "BETWEEN"
	case AEXPR_NOT_BETWEEN:
		between = "NOT BETWEEN"
	case AEXPR_BETWEEN_SYM:
		between = "BETWEEN SYMMETRIC"
	case AEXPR_NOT_BETWEEN_SYM:
		between = "NOT BETWEEN SYMMETRIC"
	}

	name, err := node.Lexpr.Deparse(Context_None)
	if err != nil {
		return nil, err
	}

	rightExpression := node.Rexpr.(List)
	out := make([]string, len(rightExpression.Items))
	for i, expr := range rightExpression.Items {
		if str, err := expr.Deparse(Context_None); err != nil {
			return nil, err
		} else {
			out[i] = *str
		}
	}

	result := fmt.Sprintf("%s %s %s", *name, between, strings.Join(out, " AND "))
	return &result, nil
}

func (node A_Expr) deparseAExprNullIf(ctx Context) (*string, error) {
	leftString, err := node.Lexpr.Deparse(Context_None)
	if err != nil {
		return nil, err
	}

	rightString, err := node.Rexpr.Deparse(Context_None)
	if err != nil {
		return nil, err
	}

	result := fmt.Sprintf("NULLIF(%s, %s)", *leftString, *rightString)
	return &result, nil
}
