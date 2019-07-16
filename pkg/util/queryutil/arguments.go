package queryutil

import (
	"fmt"
	"github.com/ahmetb/go-linq"
	"github.com/elliotcourant/noahdb/pkg/ast"
	"github.com/elliotcourant/noahdb/pkg/types"
	"reflect"
	"strconv"
)

// PlaceholderTypes relates placeholder names to their resolved type.
type PlaceholderTypes map[int]types.Type

// QueryArguments relates placeholder names to their provided query argument.
//
// A nil value represents a NULL argument.
type QueryArguments []types.Value

// GetArguments returns a distinct list of argument numbers that were found in the provided query.
func GetArguments(stmt interface{}) []int {
	args := examineArguments(stmt)
	linq.From(args).Distinct().ToSlice(&args)
	return args
}

// GetArgumentsEx will pull all of the ParamRef objects that
// it can find in a query and return them as an array.
// But if the ParamRef object has a parent object that is a
// A_Expr or a TypeCase then it will include the parent object
// as the item in the array instead (with the ParamRef as a child
// object). This function is specifically used to pull params
// from the query to try to infer their types. So if we are
// casting that param to another type or if we are comparing that
// param to another object we want to assert that objects type
// so that we can then assume the param's type.
func GetArgumentsEx(stmt interface{}) []ast.Node {
	return examineArgumentsEx(nil, stmt)
}

func ReplaceArguments(stmt interface{}, args QueryArguments) interface{} {
	return replaceArguments(stmt, 0, args)
}

func examineArguments(value interface{}) []int {
	if value == nil {
		return []int{}
	}

	if paramRef, ok := value.(ast.ParamRef); ok {
		return []int{paramRef.Number}
	}

	args := make([]int, 0)
	t := reflect.TypeOf(value)
	v := reflect.ValueOf(value)
	switch t.Kind() {
	case reflect.Ptr:
		if v.Elem().IsValid() {
			return examineArguments(v.Elem().Interface())
		}
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice:
		if v.Len() > 0 {
			for i := 0; i < v.Len(); i++ {
				args = append(args, examineArguments(v.Index(i).Interface())...)
			}
		}
	case reflect.Struct:
		for i := 0; i < t.NumField(); i++ {
			args = append(args, examineArguments(reflect.ValueOf(value).Field(i).Interface())...)
		}
	}

	return args
}

func examineArgumentsEx(parent, value interface{}) []ast.Node {
	if value == nil {
		return []ast.Node{}
	}

	if param, ok := value.(ast.ParamRef); ok {
		switch parent.(type) {
		case ast.A_Expr:
			return []ast.Node{parent.(ast.Node)}
		default:
			return []ast.Node{param}
		}
	}

	t := reflect.TypeOf(value)
	v := reflect.ValueOf(value)
	args := make([]ast.Node, 0)
	switch t.Kind() {
	case reflect.Ptr:
		if v.Elem().IsValid() {
			return examineArgumentsEx(value, v.Elem().Interface())
		}
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice:
		if v.Len() > 0 {
			for i := 0; i < v.Len(); i++ {
				args = append(args, examineArgumentsEx(value, v.Index(i).Interface())...)
			}
		}
	case reflect.Struct:
		for i := 0; i < t.NumField(); i++ {
			args = append(args, examineArgumentsEx(value, reflect.ValueOf(value).Field(i).Interface())...)
		}
	}

	return args
}

func replaceArguments(value interface{}, depth int, args QueryArguments) interface{} {
	print := func(msg string, args ...interface{}) {
		// fmt.Printf("%s%s\n", strings.Repeat("\t", depth), fmt.Sprintf(msg, args...))
	}

	if param, ok := value.(ast.ParamRef); ok {
		arg := args[param.Number-1]
		return func() ast.Node {
			switch argValue := arg.Get().(type) {
			case string:
				return ast.A_Const{
					Val: ast.String{
						Str: argValue,
					},
				}
			case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
				intv, _ := strconv.ParseUint(fmt.Sprintf("%v", argValue), 10, 64)
				return ast.A_Const{
					Val: ast.Integer{
						Ival: int64(intv),
					},
				}
			case *types.Numeric:
				floatyMcFloatyFace := float64(0.0)
				argValue.AssignTo(&floatyMcFloatyFace)
				return ast.A_Const{
					Val: ast.Float{
						Str: fmt.Sprintf("%v", floatyMcFloatyFace),
					},
				}
			case float32, float64:
				return ast.A_Const{
					Val: ast.Float{
						Str: fmt.Sprintf("%v", argValue),
					},
				}
			case bool:
				boolVal := string([]rune(fmt.Sprintf("%v", argValue))[0])
				return ast.TypeCast{
					Arg: ast.A_Const{
						Val: ast.String{
							Str: boolVal,
						},
					},
					TypeName: &ast.TypeName{
						Names: ast.List{
							Items: []ast.Node{
								ast.String{
									Str: "pg_catalog",
								},
								ast.String{
									Str: "bool",
								},
							},
						},
					},
				}
			case nil:
				return &ast.A_Const{
					Val: ast.Null{},
				}
			default:
				panic(fmt.Sprintf("unsupported type %+v", argValue))
			}
		}()
	}

	if value == nil {
		return nil
	}

	typ := reflect.TypeOf(value)
	val := reflect.ValueOf(value)

	print("[-] Parent Type <%s> Kind <%s>", typ.Name(), typ.Kind().String())
	depth++
	switch typ.Kind() {
	case reflect.Ptr:
		return value
	case reflect.Slice:
		if val.Len() > 0 {
			print("[-] Slice Type <%s> Size: %d", val.Type().String(), val.Len())
			copySlice := reflect.MakeSlice(reflect.SliceOf(typ.Elem()), val.Len(), (val.Cap()+1)*2)
			reflect.Copy(copySlice, val)
			print("[")
			depth++
			for i := 0; i < val.Len(); i++ {
				item := val.Index(i)
				copy := reflect.New(item.Type())
				print("[%d] Copying Item Type <%s> Kind <%s> Actual %s", i, copy.Type().String(), copy.Kind().String(), item.Type().String())
				result := replaceArguments(item.Interface(), depth+1, args)
				if result != nil {
					copySlice.Index(i).Set(reflect.ValueOf(result))
				}
			}
			depth--
			print("]")
			return copySlice.Interface()
		} else {
			print("[-] Slice Type <%s> Size: Empty", val.Type().Name())
		}
	case reflect.Struct:
		copy := reflect.New(typ).Elem()
		copyType := reflect.TypeOf(copy.Interface())
		print("[-] Copied Struct Type <%s> Kind <%s> Fields: %d", copy.Type().Name(), copy.Kind().String(), copy.NumField())
		depth++
		for i := 0; i < copy.NumField(); i++ {
			field := copy.Field(i)
			fieldType := copyType.Field(i)
			actual := val.Field(i)
			print("[%d] Copying Field <%s> Type <%s> Kind <%s> Actual %s", i, fieldType.Name, field.Type().String(), field.Kind().String(), actual.Type().String())
			result := replaceArguments(actual.Interface(), depth+1, args)
			if result != nil {
				copy.Field(i).Set(reflect.ValueOf(result))
			}
		}
		return copy.Interface()
	default:
		return value
	}
	return value
}
