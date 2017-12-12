#!/usr/bin/env bash

set -ax
set -e

SCRIPTS_DIR=$(dirname $0)

cd $SCRIPTS_DIR/..

go test ./converter/...
go test ./types/...
go test ./tests
go test ./imports
go test ./yaml
