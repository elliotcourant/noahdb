.PHONY: default strings protos embedded test coverage generated

CORE_DIRECTORY = ./pkg/core
PGERROR_DIRECTORY = ./pkg/pgerror
BUILD_DIRECTORY = ./bin
PACKAGE = github.com/elliotcourant/noahdb
EXECUTABLE_NAME = noah

default: dependencies test

dependencies: generated
	go get -t -v ./...

test:
	go test -v ./...

setup_build_dir:
	mkdir -p $(BUILD_DIRECTORY)

build: dependencies setup_build_dir
	go build -o $(BUILD_DIRECTORY)/$(EXECUTABLE_NAME) $(PACKAGE)

fresh: dependencies setup_build_dir
	go build -a -x -v -o $(BUILD_DIRECTORY)/$(EXECUTABLE_NAME) $(PACKAGE)

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
	@protoc -I=$(CORE_DIRECTORY) --go_out=$(CORE_DIRECTORY) $(CORE_DIRECTORY)/shard.proto
	@protoc -I=$(CORE_DIRECTORY) --go_out=$(CORE_DIRECTORY) $(CORE_DIRECTORY)/data_node.proto
	@protoc -I=$(CORE_DIRECTORY) --go_out=$(CORE_DIRECTORY) $(CORE_DIRECTORY)/tenant.proto
	@protoc -I=$(CORE_DIRECTORY) --go_out=$(CORE_DIRECTORY) $(CORE_DIRECTORY)/type.proto
	@protoc -I=$(CORE_DIRECTORY) --go_out=$(CORE_DIRECTORY) $(CORE_DIRECTORY)/table.proto
	@protoc -I=$(CORE_DIRECTORY) --go_out=$(CORE_DIRECTORY) $(CORE_DIRECTORY)/setting.proto
	@protoc -I=$(CORE_DIRECTORY) --go_out=$(CORE_DIRECTORY) $(CORE_DIRECTORY)/schema.proto
	@protoc -I=$(CORE_DIRECTORY) --go_out=$(CORE_DIRECTORY) $(CORE_DIRECTORY)/user.proto
	@protoc -I=$(PGERROR_DIRECTORY) --go_out=$(PGERROR_DIRECTORY) $(PGERROR_DIRECTORY)/errors.proto

embedded:
	@echo generating embedded files...
	@go get -u -a github.com/elliotcourant/statik
	@statik -src=$(CORE_DIRECTORY)/static/files -dest $(CORE_DIRECTORY) -f -p static

generated: strings protos embedded
