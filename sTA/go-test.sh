#!/bin/zsh

# this script will run the validation tests

go clean -testcache
go test ./internal/test/... -v
