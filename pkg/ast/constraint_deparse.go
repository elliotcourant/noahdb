// Auto-generated - DO NOT EDIT

package ast

import (
	"fmt"
	"strings"
)

func (node Constraint) Deparse(ctx Context) (string, error) {
	out := make([]string, 0)
	if node.Conname != nil {
		out = append(out, "CONSTRAINT")
		out = append(out, *node.Conname)
	}
	switch node.Contype {
	case CONSTR_NULL:
		out = append(out, "NULL")
	case CONSTR_NOTNULL:
		out = append(out, "NOT NULL")
	case CONSTR_DEFAULT:
		out = append(out, "DEFAULT")
	case CONSTR_CHECK:
		out = append(out, "CHECK")
	case CONSTR_PRIMARY:
		out = append(out, "PRIMARY KEY")
	case CONSTR_UNIQUE:
		out = append(out, "UNIQUE")
	case CONSTR_EXCLUSION:
		out = append(out, "EXCLUSION")
	case CONSTR_FOREIGN:
		if node.Conname != nil {
			out = append(out, "FOREIGN KEY")
		}
	}

	if node.RawExpr != nil {
		if expr, err := node.RawExpr.Deparse(Context_None); err != nil {
			return "", err
		} else {
			if aexpr, ok := node.RawExpr.(A_Expr); ok && aexpr.Kind == AEXPR_OP {
				out = append(out, fmt.Sprintf("(%s)", expr))
			} else {
				out = append(out, expr)
			}
		}
	}

	if node.Keys.Items != nil && len(node.Keys.Items) > 0 {
		if list, err := deparseNodeList(node.Keys.Items, Context_None); err != nil {
			return "", err
		} else {
			out = append(out, fmt.Sprintf("(%s)", strings.Join(list, ", ")))
		}
	}

	if node.FkAttrs.Items != nil && len(node.FkAttrs.Items) > 0 {
		if list, err := deparseNodeList(node.FkAttrs.Items, Context_None); err != nil {
			return "", err
		} else {
			out = append(out, fmt.Sprintf("(%s)", strings.Join(list, ", ")))
		}
	}

	if node.Pktable != nil {
		if list, err := deparseNodeList(node.PkAttrs.Items, Context_None); err != nil {
			return "", err
		} else {
			if pk, err := node.Pktable.Deparse(Context_None); err != nil {
				return "", err
			} else {
				out = append(out, fmt.Sprintf("REFERENCES %s (%s)", pk, strings.Join(list, ", ")))
			}

			switch node.FkDelAction {
			case 97: // Default (NO ACTION)
			case 99: // On Delete cascade
				out = append(out, "ON DELETE CASCADE")
			case 110:
				out = append(out, "ON DELETE SET NULL")
			case 114: // On Delete Restrict
				out = append(out, "ON DELETE RESTRICT")
			}
		}
	}

	if node.SkipValidation {
		out = append(out, "NOT VALID")
	}

	if node.Indexname != nil {
		out = append(out, fmt.Sprintf("USING INDEX %s", *node.Indexname))
	}
	return strings.Join(out, " "), nil
}
