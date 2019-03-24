// Auto-generated - DO NOT EDIT

package ast

import (
	"fmt"
	"reflect"
	"strings"
)

func (node TypeName) Deparse(ctx Context) (string, error) {
	if node.Names.Items == nil || len(node.Names.Items) == 0 {
		return "", fmt.Errorf("cannot have no names on type name")
	}
	names := make([]string, len(node.Names.Items))
	for i, name := range node.Names.Items {
		if str, err := name.Deparse(Context_TypeName); err != nil {
			return "", err
		} else {
			names[i] = str
		}
	}

	// Intervals are tricky and should be handled in a seperate method because they require some bitmask operations
	if reflect.DeepEqual(names, []string{"pg_catalog", "interval"}) {
		return node.deparseIntervalType()
	}

	out := make([]string, 0)
	if node.Setof {
		out = append(out, "SETOF")
	}

	args := ""
	if node.Typmods.Items != nil && len(node.Typmods.Items) > 0 {
		arguments := make([]string, len(node.Typmods.Items))
		for i, arg := range node.Typmods.Items {
			if str, err := arg.Deparse(Context_None); err != nil {
				return "", err
			} else {
				arguments[i] = str
			}
		}
		args = strings.Join(arguments, ", ")
	}

	if str, err := node.deparseTypeNameCase(names, args); err != nil {
		return "", err
	} else {
		out = append(out, str)
	}

	if node.ArrayBounds.Items != nil || len(node.ArrayBounds.Items) > 0 {
		out[len(out)-1] = fmt.Sprintf("%s[]", out[len(out)-1])
	}

	return strings.Join(out, ", "), nil
}

func (node TypeName) deparseIntervalType() (string, error) {
	out := []string{"interval"}

	if node.Typmods.Items != nil && len(node.Typmods.Items) > 0 {
		return "", nil
		// In the ruby version of this code this was here to
		// handle `interval hour to second(5)` but i've not
		// ever seen that syntax and will come back to it
	}

	return strings.Join(out, " "), nil
}

func (node TypeName) deparseTypeNameCase(names []string, arguments string) (string, error) {
	if names[0] != "pg_catalog" {
		return strings.Join(names, "."), nil
	}

	switch names[len(names)-1] {
	case "bpchar":
		if len(arguments) == 0 {
			return "char", nil
		} else {
			return fmt.Sprintf("char(%s)", arguments), nil
		}
	case "varchar":
		if len(arguments) == 0 {
			return "varchar", nil
		} else {
			return fmt.Sprintf("varchar(%s)", arguments), nil
		}
	case "numeric":
		if len(arguments) == 0 {
			return "numeric", nil
		} else {
			return fmt.Sprintf("numeric(%s)", arguments), nil
		}
	case "bool":
		return "boolean", nil
	case "int2":
		return "smallint", nil
	case "int4":
		return "int", nil
	case "int8":
		return "bigint", nil
	case "real", "float4":
		return "real", nil
	case "float8":
		return "double", nil
	case "time":
		return "time", nil
	case "timezt":
		return "time with time zone", nil
	case "timestamp":
		return "timestamp", nil
	case "timestamptz":
		return "timestamp with time zone", nil
	default:
		return "", fmt.Errorf("cannot deparse type: %s", names[len(names)-1])
	}
}
