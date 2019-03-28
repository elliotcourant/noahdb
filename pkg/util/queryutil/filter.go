package queryutil

import (
	"github.com/elliotcourant/noahdb/pkg/ast"
	"reflect"
)

func FindAccountIds(stmt interface{}, shardColumnName string) []uint64 {
	return examineWhereClause(stmt, 0, shardColumnName)
}

func examineWhereClause(value interface{}, depth int, shardColumnName string) []uint64 {
	ids := make([]uint64, 0)
	print := func(msg string, args ...interface{}) {
		// fmt.Printf("%s%s\n", strings.Repeat("\t", depth), fmt.Sprintf(msg, args...))
	}

	if value == nil {
		return ids
	}

	t := reflect.TypeOf(value)
	v := reflect.ValueOf(value)

	if v.Type() == reflect.TypeOf(ast.A_Expr{}) {
		expr := value.(ast.A_Expr)
		colRef, ok := expr.Lexpr.(ast.ColumnRef)
		if !ok {
			goto Continue
		}

		if colRef.Fields.Items[len(colRef.Fields.Items)-1].(ast.String).Str != shardColumnName {
			goto Continue
		}

		valueConst, ok := expr.Rexpr.(ast.A_Const)
		if !ok {
			goto Continue
		}

		numericValue, ok := valueConst.Val.(ast.Integer)
		if !ok {
			goto Continue
		}
		ids = append(ids, uint64(numericValue.Ival))
		return ids
	}

Continue:
	switch t.Kind() {
	case reflect.Ptr:
		if v.Elem().IsValid() {
			ids = append(ids, examineWhereClause(v.Elem().Interface(), depth+1, shardColumnName)...)
		}
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice:
		depth--
		if v.Len() > 0 {
			print("[")
			for i := 0; i < v.Len(); i++ {
				depth++
				print("[%d] Type {%s} {", i, v.Index(i).Type().String())
				ids = append(ids, examineWhereClause(v.Index(i).Interface(), depth+1, shardColumnName)...)
				print("},")
				depth--
			}
			print("]")
		} else {
			print("[]")
		}
	case reflect.Struct:
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			print("[%d] Field {%s} Type {%s} Kind {%s}", i, f.Name, f.Type.String(), reflect.ValueOf(value).Field(i).Kind().String())
			ids = append(ids, examineWhereClause(reflect.ValueOf(value).Field(i).Interface(), depth+1, shardColumnName)...)
		}
	}
	return ids
}
