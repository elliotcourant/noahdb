// Auto-generated - DO NOT EDIT

package pg_query

import (
	"errors"
	"fmt"
	"strings"
)

var (
	sortDirection = map[SortByDir]string{
		SORTBY_DEFAULT: "",
		SORTBY_DESC:    "DESC",
		SORTBY_ASC:     "ASC",
	}
)

func (node SortBy) Deparse(ctx Context) (*string, error) {
	out := []string{""}

	if str, err := node.Node.Deparse(Context_None); err != nil {
		return nil, err
	} else {
		out[0] = *str
	}

	if dir, ok := sortDirection[node.SortbyDir]; !ok {
		return nil, errors.New(fmt.Sprintf("cannot handle sort direction [%s]", node.SortbyDir.String()))
	} else {
		if dir != "" {
			out = append(out, dir)
		}
	}

	result := strings.Join(out, " ")
	return &result, nil
}
