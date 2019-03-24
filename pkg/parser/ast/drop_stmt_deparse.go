// Auto-generated - DO NOT EDIT

package ast

import (
	"fmt"
	"strings"
)

var (
	dropStmtRemoveTypes = map[ObjectType]string{
		OBJECT_ACCESS_METHOD: "ACCESS METHOD",
		OBJECT_AGGREGATE:     "AGGREGATE",
		OBJECT_CAST:          "CAST",
		OBJECT_COLLATION:     "COLLATION",
		OBJECT_CONVERSION:    "CONVERSION",
		OBJECT_DATABASE:      "DATABASE", // technically this gets handled by dropdb_stmt.go

		OBJECT_TABLE: "TABLE",
	}
)

func (node DropStmt) Deparse(ctx Context) (string, error) {
	out := []string{"DROP", ""}

	if removeType, ok := dropStmtRemoveTypes[node.RemoveType]; !ok {
		panic(fmt.Sprintf("cannot handle remove type [%s]", node.RemoveType.String()))
	} else {
		out[1] = removeType
	}

	if node.MissingOk {
		out = append(out, "IF EXISTS")
	}

	objects := make([]string, len(node.Objects.Items))
	for i, obj := range node.Objects.Items {
		switch obj.(type) {
		case List:
			list := obj.(List)
			if objs, err := list.DeparseList(Context_None); err != nil {
				return "", err
			} else {
				switch node.RemoveType {
				case OBJECT_CAST:
					objects[i] = fmt.Sprintf("(%s)", strings.Join(objs, " AS "))
				default:
					objects[i] = strings.Join(objs, ".")
				}
			}
		default:
			if str, err := obj.Deparse(Context_None); err != nil {
				return "", err
			} else {
				objects[i] = str
			}
		}
	}

	out = append(out, strings.Join(objects, ", "))

	switch node.Behavior {
	case DROP_CASCADE:
		out = append(out, "CASCADE")
	case DROP_RESTRICT:
		// By default the drop will restrict, so there is no need to have this in there.
	}

	return strings.Join(out, " "), nil
}
