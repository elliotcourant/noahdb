#!/usr/bin/env bash

set -e
echo "" > coverage.txt

go test -v -race -coverprofile=profile.out ./...
cat profile.out >> coverage.txt
rm profile.out
#
#for d in $(go list ./... | grep -v vendor | grep -v pkg/ast); do
#    go test -v -race -coverprofile=profile.out -covermode=atomic "$d"
#    if [ -f profile.out ]; then
#        cat profile.out >> coverage.txt
#        rm profile.out
#    fi
#done