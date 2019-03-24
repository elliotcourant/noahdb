// Auto-generated - DO NOT EDIT

package pg_query

import (
	"github.com/juju/errors"
	"strings"
)

func (node WithClause) Deparse(ctx Context) (*string, error) {
	out := []string{"WITH"}
	if node.Recursive {
		out = append(out, "RECURSIVE")
	}

	if node.Ctes.Items == nil || len(node.Ctes.Items) == 0 {
		return nil, errors.New("cannot have with clause without ctes")
	}

	ctes := make([]string, len(node.Ctes.Items))
	for i, cte := range node.Ctes.Items {
		if str, err := deparseNode(cte, Context_None); err != nil {
			return nil, err
		} else {
			ctes[i] = *str
		}
	}
	out = append(out, strings.Join(ctes, ", "))
	result := strings.Join(out, " ")
	return &result, nil
}
