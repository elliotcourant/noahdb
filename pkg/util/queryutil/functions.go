package queryutil

import (
	"fmt"
	"github.com/ahmetb/go-linq/v3"
	"github.com/elliotcourant/noahdb/pkg/ast"
	"github.com/elliotcourant/noahdb/pkg/types"
	"reflect"
	"strconv"
	"strings"
)

type BuiltInFunction func(args ...ast.Node) (interface{}, error)
type BuiltInFunctionMap map[string]BuiltInFunction

func ReplaceFunctionCalls(stmt interface{}, builtIns BuiltInFunctionMap) (interface{}, error) {
	return replaceFunctions(stmt, 0, builtIns)
}

func replaceFunctions(value interface{}, depth int, builtIns BuiltInFunctionMap) (interface{}, error) {
	print := func(msg string, args ...interface{}) {
		// fmt.Printf("%s%s\n", strings.Repeat("\t", depth), fmt.Sprintf(msg, args...))
	}

	if funcCall, ok := value.(ast.FuncCall); ok {
		functionName, err := funcCall.Name()
		if err != nil {
			return funcCall, err
		}
		print("found function [%s] checking for drop-ins", functionName)
		if builtInFunction, ok := builtIns[functionName]; ok {
			result, err := builtInFunction(funcCall.Args.Items...)
			if err != nil {
				return funcCall, err
			}
			// names := strings.Split(functionName, ".")
			// name := names[len(names) - 1]
			val, err := convertObjectToQueryLiteral(result)
			if err != nil {
				return nil, err
			}
			return val, nil
		}
	}

	if value == nil {
		return nil, nil
	}

	typ := reflect.TypeOf(value)
	val := reflect.ValueOf(value)

	print("[-] Parent Type <%s> Kind <%s>", typ.Name(), typ.Kind().String())
	depth++
	switch typ.Kind() {
	case reflect.Ptr:
		return value, nil
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
				result, err := replaceFunctions(item.Interface(), depth+1, builtIns)
				if err != nil {
					return nil, err
				}
				if result != nil {
					copySlice.Index(i).Set(reflect.ValueOf(result))
				}
			}
			depth--
			print("]")
			return copySlice.Interface(), nil
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
			result, err := replaceFunctions(actual.Interface(), depth+1, builtIns)
			if err != nil {
				return nil, err
			}
			if result != nil {
				copy.Field(i).Set(reflect.ValueOf(result))
			}
		}
		return copy.Interface(), nil
	default:
		return value, nil
	}
	return value, nil
}

func convertObjectToQueryLiteral(obj interface{}) (ast.Node, error) {
	switch argValue := obj.(type) {
	case string:
		return ast.A_Const{
			Val: ast.String{
				Str: argValue,
			},
		}, nil
	case []string:
		items := make([]string, len(argValue))
		linq.From(argValue).SelectT(func(val string) string {
			return fmt.Sprintf(`"%s"`, strings.Replace(val, `"`, `\"`, 0))
		}).ToSlice(&items)
		return ast.TypeCast{
			Arg: ast.A_Const{
				Val: ast.String{
					Str: fmt.Sprintf("{%s}", strings.Join(items, ", ")),
				},
			},
			TypeName: &ast.TypeName{
				Names: ast.List{
					Items: []ast.Node{
						ast.String{
							Str: "text",
						},
					},
				},
				ArrayBounds: ast.List{
					Items: []ast.Node{
						ast.Integer{
							Ival: -1,
						},
					},
				},
			},
		}, nil
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		intv, _ := strconv.ParseUint(fmt.Sprintf("%v", argValue), 10, 64)
		return ast.A_Const{
			Val: ast.Integer{
				Ival: int64(intv),
			},
		}, nil
	case *types.Numeric:
		floatyMcFloatyFace := float64(0.0)
		err := argValue.AssignTo(&floatyMcFloatyFace)
		if err != nil {
			return nil, err
		}
		return ast.A_Const{
			Val: ast.Float{
				Str: fmt.Sprintf("%v", floatyMcFloatyFace),
			},
		}, nil
	case float32, float64:
		return ast.A_Const{
			Val: ast.Float{
				Str: fmt.Sprintf("%v", argValue),
			},
		}, nil
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
		}, nil
	case nil:
		return &ast.A_Const{
			Val: ast.Null{},
		}, nil
	default:
		panic(fmt.Sprintf("unsupported type %+v", argValue))
	}
}
