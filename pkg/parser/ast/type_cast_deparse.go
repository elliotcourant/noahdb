// Auto-generated - DO NOT EDIT

package pg_query

import (
	"fmt"
	"github.com/juju/errors"
)

func (node TypeCast) Deparse(ctx Context) (*string, error) {
	if node.TypeName == nil {
		return nil, errors.New("typename cannot be null in typecast")
	}
	if str, err := deparseNode(*node.TypeName, Context_None); err != nil {
		return nil, err
	} else {
		if val, err := deparseNode(node.Arg, Context_None); err != nil {
			return nil, err
		} else {
			if *str == "boolean" {
				if *val == "'t'" {
					result := "true"
					return &result, nil
				} else {
					result := "false"
					return &result, nil
				}
			}

			if typename, err := (*node.TypeName).Deparse(Context_None); err != nil {
				return nil, err
			} else {
				result := fmt.Sprintf("%s::%s", *val, *typename)
				return &result, nil
			}
		}
	}
}
