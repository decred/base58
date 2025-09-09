#!/usr/bin/env bash

set -ex

# This script runs the tests for all packages in the repository and then uses
# golangci-lint (github.com/golangci/golangci-lint) to run all linters defined
# by the configuration in .golangci.yml.

# run tests
env GORACE="halt_on_error=1" go test -race ./...

# run linters
golangci-lint run
