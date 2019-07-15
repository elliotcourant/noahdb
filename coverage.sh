#!/usr/bin/env bash

set -e
echo "" > coverage.txt

go test -v -race -coverprofile=coverage.txt -coverpkg=./... -covermode=atomic ./...