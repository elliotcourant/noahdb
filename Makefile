a.PHONY: default strings protos embedded test coverage generated docker

PKG_DIRECTORY := pkg
CORE_DIRECTORY := pkg/core
PGERROR_DIRECTORY := pkg/pgerror
TYPES_DIRECTORY := pkg/types
BUILD_DIRECTORY := bin
PACKAGE = github.com/elliotcourant/noahdb
EXECUTABLE_NAME = noah
DOCKER_TAG = local

postgres:
	docker build -t noahdb/postgres:local ./k8s/postgres

docker:
	docker build -t noahdb/node:$(DOCKER_TAG) -f ./k8s/node/Dockerfile .

docker-test:
	docker build -t noahdb/test:local -f ./k8s/node/Dockerfile.test .

kube: docker
	kubectl delete -f ./k8s/node/noahdb.yaml --wait --ignore-not-found=true
	sleep 5
	kubectl apply -f ./k8s/node/noahdb.yaml

default: generated test

test:
	go test -race -v ./...

setup_build_dir:
	mkdir -p $(BUILD_DIRECTORY)

build: generated setup_build_dir
	go build -o $(BUILD_DIRECTORY)/$(EXECUTABLE_NAME) $(PACKAGE)

fresh: generated setup_build_dir
	go build -a -x -v -o $(BUILD_DIRECTORY)/$(EXECUTABLE_NAME) $(PACKAGE)

coverage:
	./coverage.sh

strings:
	@echo generating strings...
	@go get -u -a golang.org/x/tools/cmd/stringer

	@stringer -type commandType -output pkg/frunk/command.string.go pkg/frunk/command.go
	@stringer -type ClusterState,ConsistencyLevel,BackupFormat -output pkg/frunk/store.string.go pkg/frunk/store.go

	@stringer -type Context -output pkg/ast/context.string.go pkg/ast/context.go
	@stringer -type ConstrType -output pkg/ast/constr_type.string.go pkg/ast/constr_type.go
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
	@protoc -I=$(PKG_DIRECTORY) --go_out=$(PKG_DIRECTORY) $(CORE_DIRECTORY)/table.proto
	@protoc -I=$(CORE_DIRECTORY) --go_out=$(CORE_DIRECTORY) $(CORE_DIRECTORY)/setting.proto
	@protoc -I=$(CORE_DIRECTORY) --go_out=$(CORE_DIRECTORY) $(CORE_DIRECTORY)/schema.proto
	@protoc -I=$(CORE_DIRECTORY) --go_out=$(CORE_DIRECTORY) $(CORE_DIRECTORY)/user.proto
	@protoc -I=$(PGERROR_DIRECTORY) --go_out=$(PGERROR_DIRECTORY) $(PGERROR_DIRECTORY)/errors.proto
	@protoc -I=$(TYPES_DIRECTORY) --go_out=$(TYPES_DIRECTORY) $(TYPES_DIRECTORY)/type.proto

embedded:
	@echo generating embedded files...
	@go get -u -a github.com/elliotcourant/statik@master
	@statik -src=$(CORE_DIRECTORY)/static/files -dest $(CORE_DIRECTORY) -f -p static

generated: strings protos embedded
