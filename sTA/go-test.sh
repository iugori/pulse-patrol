#!/bin/zsh

go clean -testcache
go test ./internal/test/... -v
