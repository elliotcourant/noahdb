.PHONY: default

DOCKER_TAG = local

default: docker

docker:
		docker build -t noahdb/postgres:$(DOCKER_TAG) .
