.PHONY: default strings protos

default: strings

strings:
	@go get -u -a golang.org/x/tools/cmd/stringer
	@stringer -type Context -output pkg/ast/context.string.go pkg/ast/context.go
	@stringer -type ObjectType -output pkg/ast/object_type.string.go pkg/ast/object_type.go
	@stringer -type SortByDir -output pkg/ast/sort_by_dir.string.go pkg/ast/sort_by_dir.go
	@stringer -type StmtType -output pkg/ast/stmt_type.string.go pkg/ast/stmt_type.go
	@stringer -type SubLinkType -output pkg/ast/sub_link_type.string.go pkg/ast/sub_link_type.go
	@stringer -type SQLValueFunctionOp -output pkg/ast/sql_value_function_op.string.go pkg/ast/sql_value_function_op.go
	@stringer -type TransactionStmtKind -output pkg/ast/transaction_stmt_kind.string.go pkg/ast/transaction_stmt_kind.go

STORE_DIRECTORY = ./pkg/store

protos:
	protoc -I=$(STORE_DIRECTORY) --go_out=$(STORE_DIRECTORY) $(STORE_DIRECTORY)/action.proto
	protoc -I=$(STORE_DIRECTORY) --go_out=$(STORE_DIRECTORY) $(STORE_DIRECTORY)/sequence.proto