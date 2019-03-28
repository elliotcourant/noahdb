package queryutil

import (
	"github.com/elliotcourant/noahdb/pkg/ast"
	"reflect"
)

func GetColumns(stmt ast.Stmt) []ast.ResTarget {
	cols := make([]ast.ResTarget, 0)
	switch tree := stmt.(type) {
	case ast.SelectStmt:
		return examineColumns(tree.TargetList, 0)
	case ast.InsertStmt:
		return examineColumns(tree.ReturningList, 0)
	case ast.UpdateStmt:
		return examineColumns(tree.ReturningList, 0)
	case ast.DeleteStmt:
		return examineColumns(tree.ReturningList, 0)
	}
	return cols
}

func examineColumns(value interface{}, depth int) []ast.ResTarget {
	cols := make([]ast.ResTarget, 0)
	print := func(msg string, args ...interface{}) {
		// fmt.Printf("%s%s\n", strings.Repeat("\t", depth), fmt.Sprintf(msg, args...))
	}

	if value == nil {
		return cols
	}

	t := reflect.TypeOf(value)
	v := reflect.ValueOf(value)

	if v.Type() == reflect.TypeOf(ast.ResTarget{}) {
		col := value.(ast.ResTarget)
		cols = append(cols, col)
	}

	switch t.Kind() {
	case reflect.Ptr:
		if v.Elem().IsValid() {
			cols = append(cols, examineColumns(v.Elem().Interface(), depth+1)...)
		}
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice:
		depth--
		if v.Len() > 0 {
			print("[")
			for i := 0; i < v.Len(); i++ {
				depth++
				print("[%d] Type {%s} {", i, v.Index(i).Type().String())
				cols = append(cols, examineColumns(v.Index(i).Interface(), depth+1)...)
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
			cols = append(cols, examineColumns(reflect.ValueOf(value).Field(i).Interface(), depth+1)...)
		}
	}
	return cols
}
