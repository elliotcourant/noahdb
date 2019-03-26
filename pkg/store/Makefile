.PHONY: default build protos
default: protos

protos:
	protoc -I=./ --go_out=plugins=grpc:./ ./raft.proto
	protoc -I=./ --go_out=./ ./fsm.proto