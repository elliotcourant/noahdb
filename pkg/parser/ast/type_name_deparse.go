// Auto-generated - DO NOT EDIT

package pg_query

import (
	"fmt"
	"github.com/juju/errors"
	"reflect"
	"strings"
)

func (node TypeName) Deparse(ctx Context) (*string, error) {
	if node.Names.Items == nil || len(node.Names.Items) == 0 {
		return nil, errors.New("cannot have no names on type name")
	}
	names := make([]string, len(node.Names.Items))
	for i, name := range node.Names.Items {
		if str, err := deparseNode(name, Context_TypeName); err != nil {
			return nil, err
		} else {
			names[i] = *str
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
			if str, err := deparseNode(arg, Context_None); err != nil {
				return nil, err
			} else {
				arguments[i] = *str
			}
		}
		args = strings.Join(arguments, ", ")
	}

	if str, err := node.deparseTypeNameCase(names, args); err != nil {
		return nil, err
	} else {
		out = append(out, *str)
	}

	if node.ArrayBounds.Items != nil || len(node.ArrayBounds.Items) > 0 {
		out[len(out)-1] = fmt.Sprintf("%s[]", out[len(out)-1])
	}

	result := strings.Join(out, ", ")
	return &result, nil
}

func (node TypeName) deparseIntervalType() (*string, error) {
	out := []string{"interval"}

	if node.Typmods.Items != nil && len(node.Typmods.Items) > 0 {
		return nil, nil
		// In the ruby version of this code this was here to
		// handle `interval hour to second(5)` but i've not
		// ever seen that syntax and will come back to it
	}

	result := strings.Join(out, " ")
	return &result, nil
}

func (node TypeName) deparseTypeNameCase(names []string, arguments string) (*string, error) {
	if names[0] != "pg_catalog" {
		result := strings.Join(names, ".")
		return &result, nil
	}

	switch names[len(names)-1] {
	case "bpchar":
		if len(arguments) == 0 {
			result := "char"
			return &result, nil
		} else {
			result := fmt.Sprintf("char(%s)", arguments)
			return &result, nil
		}
	case "varchar":
		if len(arguments) == 0 {
			result := "varchar"
			return &result, nil
		} else {
			result := fmt.Sprintf("varchar(%s)", arguments)
			return &result, nil
		}
	case "numeric":
		if len(arguments) == 0 {
			result := "numeric"
			return &result, nil
		} else {
			result := fmt.Sprintf("numeric(%s)", arguments)
			return &result, nil
		}
	case "bool":
		result := "boolean"
		return &result, nil
	case "int2":
		result := "smallint"
		return &result, nil
	case "int4":
		result := "int"
		return &result, nil
	case "int8":
		result := "bigint"
		return &result, nil
	case "real", "float4":
		result := "real"
		return &result, nil
	case "float8":
		result := "double"
		return &result, nil
	case "time":
		result := "time"
		return &result, nil
	case "timezt":
		result := "time with time zone"
		return &result, nil
	case "timestamp":
		result := "timestamp"
		return &result, nil
	case "timestamptz":
		result := "timestamp with time zone"
		return &result, nil
	default:
		return nil, errors.Errorf("cannot deparse type: %s", names[len(names)-1])
	}
	return nil, nil
}
