#!/bin/bash

set -e

src_dir=$(dirname $0)

#import version details
source $src_dir/version.sh

#build from parent dir
cd $src_dir/..

#create output dir if none exists
mkdir -p bin

#using a linkable binary for plugin support
CGO_ENABLED=1 go build -ldflags "-X github.com/koki/short/cmd.GITCOMMIT=$VERSION" -o bin/short
