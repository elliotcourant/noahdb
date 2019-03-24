// Auto-generated - DO NOT EDIT

package ast

import (
	"fmt"
)

func (node TypeCast) Deparse(ctx Context) (string, error) {
	if node.TypeName == nil {
		return "", fmt.Errorf("typename cannot be null in typecast")
	}
	if str, err := (*node.TypeName).Deparse(Context_None); err != nil {
		return "", err
	} else {
		if val, err := (node.Arg).Deparse(Context_None); err != nil {
			return "", err
		} else {
			if str == "boolean" {
				if val == "'t'" {
					return "true", nil
				} else {
					return "false", nil
				}
			}

			if typename, err := (*node.TypeName).Deparse(Context_None); err != nil {
				return "", err
			} else {
				return fmt.Sprintf("%s::%s", val, typename), nil
			}
		}
	}
}
