// Auto-generated - DO NOT EDIT

package ast

import (
	"fmt"
	"strings"
)

func (node CreateForeignServerStmt) Deparse(ctx Context) (string, error) {
	out := []string{"CREATE SERVER", ""}

	if node.Servername == nil {
		panic("server name cannot be nil in create server")
	}
	out[1] = *node.Servername

	if node.Servertype != nil {
		out = append(out, fmt.Sprintf("TYPE '%s'", *node.Servertype))
	}

	if node.Version != nil {
		out = append(out, fmt.Sprintf("VERSION '%s'", *node.Version))
	}

	if node.Fdwname == nil {
		panic("fdwname cannot be nil in create server")
	}

	out = append(out, fmt.Sprintf("FOREIGN DATA WRAPPER %s", *node.Fdwname))

	if node.Options.Items != nil && len(node.Options.Items) > 0 {
		out = append(out, "OPTIONS")
		options := make([]string, len(node.Options.Items))
		for i, option := range node.Options.Items {
			if str, err := option.Deparse(Context_None); err != nil {
				return "", err
			} else {
				options[i] = str
			}
		}
		out = append(out, fmt.Sprintf("(%s)", strings.Join(options, ", ")))
	}

	return strings.Join(out, " "), nil
}
