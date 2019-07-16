# Query Utilities

This entire package is here just to make pulling
generic information out of a query syntax tree
without making your head explode. Instead it just makes
your head hurt by using reflection!

---

## GetArgumentsEx
GetArgumentsEx will pull all of the ParamRef objects that
it can find in a query and return them as an array.
But if the ParamRef object has a parent object that is a
A_Expr or a TypeCase then it will include the parent object
as the item in the array instead (with the ParamRef as a child
object). This function is specifically used to pull params
from the query to try to infer their types. So if we are
casting that param to another type or if we are comparing that
param to another object we want to assert that objects type
so that we can then assume the param's type.

## GetArguments
GetArguments returns a distinct list of argument numbers that 
were found in the provided query.

# Benchmarks
```bash
BenchmarkGetArguments/typical_query-8         	  200000	      7471 ns/op
BenchmarkGetArguments/dead_simple_query-8     	 1000000	      1942 ns/op
BenchmarkGetArgumentsEx/typical_query-8       	  200000	      7376 ns/op
BenchmarkGetArgumentsEx/dead_simple_query-8   	 1000000	      1301 ns/op
```
