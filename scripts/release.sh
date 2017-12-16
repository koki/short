#!/usr/bin/env bash

set -e

script_dir=$(dirname $0)

#import version details
source $script_dir/version.sh

#build from parent dir
cd $script_dir/..

#create output dir if none exists
mkdir -p bin

docker run -v $GOPATH:/go -t kokster/gox gox -ldflags "-X github.com/koki/short/cmd.GITCOMMIT=$VERSION" -output="src/github.com/koki/short/bin/short_{{.OS}}_{{.Arch}}" -os="linux darwin" -arch="386 amd64" github.com/koki/short
