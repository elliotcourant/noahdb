package queryutil

import (
	"fmt"
	"github.com/ahmetb/go-linq"
	"github.com/elliotcourant/noahdb/pkg/ast"
	"reflect"
)

func FindAccountIdsEx(
	stmt interface{},
	shardColumnNames map[string]string,
	columnsAndTables map[string][]string,
) ([]uint64, error) {
	f := &findAccounts{
		shardColumnNames: shardColumnNames,
		aliases:          map[string]string{},
		columnsAndTables: columnsAndTables,
	}
	return f.findAccountIdsEx(stmt)
}

func FindAccountIds(stmt interface{}, shardColumnName string) []uint64 {
	return examineWhereClause(stmt, 0, shardColumnName)
}

type findAccounts struct {
	aliases          map[string]string
	shardColumnNames map[string]string
	columnsAndTables map[string][]string
}

func (f *findAccounts) findAccountIdsEx(value interface{}) ([]uint64, error) {
	if value == nil {
		return nil, nil
	}

	if table, ok := value.(ast.RangeVar); ok {
		aliasName, tableName := *table.Relname, *table.Relname
		if table.Alias != nil {
			aliasName = *table.Alias.Aliasname
		}
		f.aliases[aliasName] = tableName
		f.aliases[tableName] = tableName
	}

	if expr, ok := value.(ast.A_Expr); ok {
		var tableName string
		var columnName string

		var parts []string
		linq.From(expr.Lexpr.(ast.ColumnRef).Fields.Items).Select(func(i interface{}) interface{} {
			return i.(ast.String).Str
		}).ToSlice(&parts)
		// The column could belong to any table in the query.
		if len(parts) == 1 {
			columnName = parts[0]
			// The column belongs to a single table.
			if tables, ok := f.columnsAndTables[columnName]; ok && len(tables) == 1 {
				tableName = tables[0]
			} else if ok && len(tables) > 1 {
				return nil, fmt.Errorf("column [%s] is ambigious in query", columnName)
			} else {
				return nil, fmt.Errorf("column [%s] could not be resolved", columnName)
			}
		} else if len(parts) == 2 {
			tableAlias := parts[0]
			columnName = parts[1]
			if table, ok := f.aliases[tableAlias]; ok {
				tableName = table
			} else {
				return nil, fmt.Errorf("table/alias [%s] could not be resolved", tableAlias)
			}
		}

		if shardColumn, ok := f.shardColumnNames[tableName]; !ok || shardColumn != columnName {
			return nil, nil
		}

		return []uint64{uint64(expr.Rexpr.(ast.A_Const).Val.(ast.Integer).Ival)}, nil
	}

	t := reflect.TypeOf(value)
	v := reflect.ValueOf(value)
	args := make([]uint64, 0)
	switch t.Kind() {
	case reflect.Ptr:
		if v.Elem().IsValid() {
			return f.findAccountIdsEx(v.Elem().Interface())
		}
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice:
		if v.Len() > 0 {
			for i := 0; i < v.Len(); i++ {
				ids, err := f.findAccountIdsEx(v.Index(i).Interface())
				if err != nil {
					return nil, err
				}
				args = append(args, ids...)
			}
		}
	case reflect.Struct:
		for i := 0; i < t.NumField(); i++ {
			ids, err := f.findAccountIdsEx(reflect.ValueOf(value).Field(i).Interface())
			if err != nil {
				return nil, err
			}
			args = append(args, ids...)
		}
	}

	return args, nil
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
