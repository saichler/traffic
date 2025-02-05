#!/usr/bin/env bash

# Fail on errors and don't open cover file
set -e
# clean up
rm -rf go.sum
rm -rf go.mod
rm -rf vendor

# fetch dependencies
go mod init
GOPROXY=direct go mod tidy
go mod vendor

# Run unit tests with coverage
go test -v -coverpkg=./generator/... -coverprofile=cover-report.html ./... --failfast

# Open the coverage report in a browser
go tool cover -html=cover-report.html
