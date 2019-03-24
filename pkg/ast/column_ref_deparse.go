// Auto-generated - DO NOT EDIT

package ast

import (
	"github.com/juju/errors"
	"strings"
)

func (node ColumnRef) Deparse(ctx Context) (string, error) {
	if node.Fields.Items == nil || len(node.Fields.Items) == 0 {
		return "", errors.New("columnref cannot have 0 fields")
	}
	out := make([]string, len(node.Fields.Items))
	for i, field := range node.Fields.Items {
		if str, err := field.Deparse(ctx); err != nil {
			return "", err
		} else {
			out[i] = str
		}
	}
	return strings.Join(out, "."), nil
}
