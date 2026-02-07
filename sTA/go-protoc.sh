#!/bin/zsh

# go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
# go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

GOPATH_BIN="$(go env GOPATH)/bin"
export PATH="$PATH:$GOPATH_BIN"

protoc --go_out=. --go_opt=module=sTA \
       --go-grpc_out=. --go-grpc_opt=module=sTA \
       api/greeting.proto