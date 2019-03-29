.PHONY: default strings protos embedded test coverage generated

STORE_DIRECTORY = ./pkg/store
CORE_DIRECTORY = ./pkg/core
PGERROR_DIRECTORY = ./pkg/pgerror

default: test

dependencies: generated
	go get -t -v ./...

test:
	go test -v ./...

coverage:
	./coverage.sh

strings:
	@echo generating strings...
	@go get -u -a golang.org/x/tools/cmd/stringer
	@stringer -type Context -output pkg/ast/context.string.go pkg/ast/context.go
	@stringer -type ObjectType -output pkg/ast/object_type.string.go pkg/ast/object_type.go
	@stringer -type SortByDir -output pkg/ast/sort_by_dir.string.go pkg/ast/sort_by_dir.go
	@stringer -type StmtType -output pkg/ast/stmt_type.string.go pkg/ast/stmt_type.go
	@stringer -type SubLinkType -output pkg/ast/sub_link_type.string.go pkg/ast/sub_link_type.go
	@stringer -type SQLValueFunctionOp -output pkg/ast/sql_value_function_op.string.go pkg/ast/sql_value_function_op.go
	@stringer -type TransactionStmtKind -output pkg/ast/transaction_stmt_kind.string.go pkg/ast/transaction_stmt_kind.go

protos:
	@echo generating protos...
	@protoc -I=$(STORE_DIRECTORY) --go_out=plugins=grpc:$(STORE_DIRECTORY) $(STORE_DIRECTORY)/fsm.proto
	@protoc -I=$(STORE_DIRECTORY) --go_out=plugins=grpc:$(STORE_DIRECTORY) $(STORE_DIRECTORY)/raft.proto
	@protoc -I=$(CORE_DIRECTORY) --go_out=$(CORE_DIRECTORY) $(CORE_DIRECTORY)/shard.proto
	@protoc -I=$(CORE_DIRECTORY) --go_out=$(CORE_DIRECTORY) $(CORE_DIRECTORY)/data_node.proto
	@protoc -I=$(PGERROR_DIRECTORY) --go_out=$(PGERROR_DIRECTORY) $(PGERROR_DIRECTORY)/errors.proto

embedded:
	@echo generating embedded files...
	@go get -u -a github.com/elliotcourant/statik
	@statik -src=$(CORE_DIRECTORY)/static/files -dest $(CORE_DIRECTORY) -f -p static

generated: strings protos embedded