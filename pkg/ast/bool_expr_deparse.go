// Auto-generated - DO NOT EDIT

package ast

import (
	"fmt"
	"github.com/juju/errors"
	"reflect"
	"strings"
)

func (node BoolExpr) Deparse(ctx Context) (string, error) {
	// There is no BOOL_EXPR_NOT in go for some reason?
	switch node.Boolop {
	case AND_EXPR:
		return node.deparseBoolExprAnd()
	case OR_EXPR:
		return node.deparseBoolExprOr()
	case 2: // WHERE NOT
		return node.deparseBoolExprNot()
	default:
		return "", errors.Errorf("cannot handle bool expression type (%d)", node.Boolop)
	}
}

func (node BoolExpr) deparseBoolExprNot() (string, error) {
	out := []string{"NOT"}

	items := make([]string, len(node.Args.Items))
	for i, item := range node.Args.Items {
		if str, err := item.Deparse(Context_Operator); err != nil {
			return "", err
		} else {
			items[i] = str
		}
	}

	if len(items) > 1 {
		panic("cannot handle more than 1 arg in `not` expression")
	}

	out = append(out, items...)

	return strings.Join(out, " "), nil
}

func (node BoolExpr) deparseBoolExprAnd() (string, error) {
	if node.Args.Items == nil || len(node.Args.Items) == 0 {
		return "", errors.New("args cannot be empty for boolean expression")
	}
	args := make([]string, len(node.Args.Items))
	for i, arg := range node.Args.Items {
		if str, err := arg.Deparse(Context_None); err != nil {
			return "", err
		} else {
			t := reflect.TypeOf(arg)
			if t == reflect.TypeOf(BoolExpr{}) && arg.(BoolExpr).Boolop == OR_EXPR {
				args[i] = fmt.Sprintf("(%s)", str)
			} else {
				args[i] = str
			}
		}
	}
	return strings.Join(args, " AND "), nil
}

func (node BoolExpr) deparseBoolExprOr() (string, error) {
	if node.Args.Items == nil || len(node.Args.Items) == 0 {
		return "", errors.New("args cannot be empty for boolean expression")
	}
	args := make([]string, len(node.Args.Items))
	for i, arg := range node.Args.Items {
		if str, err := arg.Deparse(Context_None); err != nil {
			return "", err
		} else {
			t := reflect.TypeOf(arg)
			if t == reflect.TypeOf(BoolExpr{}) && (arg.(BoolExpr).Boolop == OR_EXPR || arg.(BoolExpr).Boolop == AND_EXPR) {
				args[i] = fmt.Sprintf("(%s)", str)
			} else {
				args[i] = str
			}
		}
	}
	return strings.Join(args, " OR "), nil
}
