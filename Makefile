.PHONY: default strings

default: strings


strings:
	@go get -u -a golang.org/x/tools/cmd/stringer
	@stringer -type Context pkg/parser/ast/context.go
	@stringer -type ObjectType pkg/parser/ast/object_type.go
	@stringer -type SortByDir pkg/parser/ast/sort_by_dir.go
	@stringer -type StmtType pkg/parser/ast/stmt_type.go
	@stringer -type SubLinkType pkg/parser/ast/sub_link_type.go
	@stringer -type SQLValueFunctionOp pkg/parser/ast/sql_value_function_op.go