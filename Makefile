.PHONY: default strings protos embedded test coverage generated docker

PKG_DIRECTORY := pkg
CORE_DIRECTORY := pkg/core
PGERROR_DIRECTORY := pkg/pgerror
TYPES_DIRECTORY := pkg/types
BUILD_DIRECTORY := bin
PACKAGE = github.com/elliotcourant/noahdb
EXECUTABLE_NAME = noah
DOCKER_TAG = development

postgres:
	docker build -t noahdb-postgres:latest ./k8s/postgres

docker:
	docker build -t noahdb:latest -f ./k8s/node/Dockerfile .

kube_down:
	kubectl delete -f ./k8s/node/noahdb.yaml --wait --ignore-not-found=true

kube: docker kube_down
	sleep 5
	kubectl apply -f ./k8s/node/noahdb.yaml

default: generated test

test: clean generated
	go test -race -v ./...

setup_build_dir:
	mkdir -p $(BUILD_DIRECTORY)

build: clean generated setup_build_dir
	go build -o $(BUILD_DIRECTORY)/$(EXECUTABLE_NAME) $(PACKAGE)

fresh: clean generated setup_build_dir
	go build -a -x -v -o $(BUILD_DIRECTORY)/$(EXECUTABLE_NAME) $(PACKAGE)

coverage:
	./coverage.sh

clean:
	rm -rfd bin
	rm -rfd vendor
	rm -rf pkg/ast/*.string.go
	rm -rf $(CORE_DIRECTORY)/*.pb.go
	rm -rf $(PGERROR_DIRECTORY)/*.pb.go
	rm -rf $(TYPES_DIRECTORY)/*.pb.go

strings:
	@echo generating strings...
	@go get -u -a golang.org/x/tools/cmd/stringer

	cd pkg/plan && make strings

	@stringer -type commandType -output pkg/frunk/command.string.go pkg/frunk/command.go
	@stringer -type ClusterState,ConsistencyLevel,BackupFormat -output pkg/frunk/store.string.go pkg/frunk/store.go

	@stringer -type TransactionStatus -output pkg/pgproto/transaction_status.string.go pkg/pgproto/transaction_status.go

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
	@protoc -I=$(TYPES_DIRECTORY) --go_out=${GOPATH}/src $(TYPES_DIRECTORY)/type.proto
	@protoc -I=$(PGERROR_DIRECTORY) --go_out=$(PGERROR_DIRECTORY) $(PGERROR_DIRECTORY)/errors.proto

embedded:
	@echo generating embedded files...
	@go get -u -a github.com/elliotcourant/statik@master
	@statik -src=$(CORE_DIRECTORY)/static/files -dest $(CORE_DIRECTORY) -f -p static

generated: strings protos embedded
