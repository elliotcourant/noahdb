// Auto-generated - DO NOT EDIT

package pg_query

import (
	"github.com/juju/errors"
	"strings"
)

func (node NullTest) Deparse(ctx Context) (*string, error) {
	out := make([]string, 0)
	if node.Arg == nil {
		return nil, errors.New("argument cannot be null for null test (ironically)")
	}

	if str, err := deparseNode(node.Arg, Context_None); err != nil {
		return nil, err
	} else {
		out = append(out, *str)
	}

	switch node.Nulltesttype {
	case IS_NULL:
		out = append(out, "IS NULL")
	case IS_NOT_NULL:
		out = append(out, "IS NOT NULL")
	default:
		return nil, errors.Errorf("could not parse null test type (%d)", node.Nulltesttype)
	}

	result := strings.Join(out, " ")
	return &result, nil
}
