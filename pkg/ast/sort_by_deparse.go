// Auto-generated - DO NOT EDIT

package ast

import (
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

func (node SortBy) Deparse(ctx Context) (string, error) {
	out := []string{""}

	if str, err := node.Node.Deparse(Context_None); err != nil {
		return "", err
	} else {
		out[0] = str
	}

	if dir, ok := sortDirection[node.SortbyDir]; !ok {
		return "", fmt.Errorf("cannot handle sort direction [%s]", node.SortbyDir.String())
	} else if dir != "" {
		out = append(out, dir)
	}

	return strings.Join(out, " "), nil
}
