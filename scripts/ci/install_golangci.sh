#!/usr/bin/env bash

set -ex

if ! which golangci-lint &> /dev/null
then
  echo 'installing golangci-lint-latest'
  curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ${GOBIN} latest
  echo 'Successfully installed'
fi
