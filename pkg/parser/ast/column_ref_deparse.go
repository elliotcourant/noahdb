// Auto-generated - DO NOT EDIT

package pg_query

import (
	"github.com/juju/errors"
	"strings"
)

func (node ColumnRef) Deparse(ctx Context) (*string, error) {
	if node.Fields.Items == nil || len(node.Fields.Items) == 0 {
		return nil, errors.New("columnref cannot have 0 fields")
	}
	out := make([]string, len(node.Fields.Items))
	for i, field := range node.Fields.Items {
		if str, err := field.Deparse(ctx); err != nil {
			return nil, err
		} else {
			out[i] = *str
		}
	}
	result := strings.Join(out, ".")
	return &result, nil
}
