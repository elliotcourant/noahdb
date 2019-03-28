package queryutil

import (
	"fmt"
	"github.com/elliotcourant/noahdb/pkg/ast"
	"github.com/elliotcourant/noahdb/pkg/types"
	"github.com/readystock/golinq"
	"reflect"
	"strconv"
)

// PlaceholderTypes relates placeholder names to their resolved type.
type PlaceholderTypes map[string]types.T

// QueryArguments relates placeholder names to their provided query argument.
//
// A nil value represents a NULL argument.
type QueryArguments map[string]types.Value

func GetArguments(stmt interface{}) []int {
	args := make([]int, 0)
	args = append(args, examineArguments(stmt, 0)...)
	linq.From(args).Distinct().ToSlice(&args)
	return args
}

func ReplaceArguments(stmt interface{}, args QueryArguments) interface{} {
	return replaceArguments(stmt, 0, args)
}

func examineArguments(value interface{}, depth int) []int {
	args := make([]int, 0)
	print := func(msg string, args ...interface{}) {
		// fmt.Printf("%s%s\n", strings.Repeat("\t", depth), fmt.Sprintf(msg, args...))
	}

	if value == nil {
		return args
	}

	t := reflect.TypeOf(value)
	v := reflect.ValueOf(value)

	if v.Type() == reflect.TypeOf(ast.ParamRef{}) {
		param := value.(ast.ParamRef)
		args = append(args, param.Number)
	}

	switch t.Kind() {
	case reflect.Ptr:
		if v.Elem().IsValid() {
			args = append(args, examineArguments(v.Elem().Interface(), depth+1)...)
		}
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice:
		depth--
		if v.Len() > 0 {
			print("[")
			for i := 0; i < v.Len(); i++ {
				depth++
				print("[%d] Type {%s} {", i, v.Index(i).Type().String())
				args = append(args, examineArguments(v.Index(i).Interface(), depth+1)...)
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
			args = append(args, examineArguments(reflect.ValueOf(value).Field(i).Interface(), depth+1)...)
		}
	}
	return args
}

func replaceArguments(value interface{}, depth int, args QueryArguments) interface{} {
	print := func(msg string, args ...interface{}) {
		// fmt.Printf("%s%s\n", strings.Repeat("\t", depth), fmt.Sprintf(msg, args...))
	}

	if param, ok := value.(ast.ParamRef); ok {
		if arg, ok := args[strconv.FormatInt(int64(param.Number), 10)]; ok {
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
		} else {
			panic("parameter is not a param reference")
		}

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
