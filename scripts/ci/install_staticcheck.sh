#!/usr/bin/env bash

set -ex

if ! which staticcheck &> /dev/null
then
  echo 'installing honnef.co/go/tools/cmd/staticcheck...'
  GO111MODULE=off go get honnef.co/go/tools/cmd/staticcheck 
fi
