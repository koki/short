#!/usr/bin/env bash

set -e

script_dir=$(dirname $0)

#import version details
source $script_dir/build.sh

#build from parent dir
cd $script_dir/..

#create output dir if none exists
mkdir -p bin

GOOS=linux GOARCH=amd64 ./scripts/build.sh 
mv bin/short bin/short_linux_amd64

GOOS=linux GOARCH=386 ./scripts/build.sh
mv bin/short bin/short_linux_386

GOOS=darwin GOARCH=amd64 ./scripts/build.sh
mv bin/short bin/short_darwin_amd64

GOOS=darwin GOARCH=386 ./scripts/build.sh
mv bin/short bin/short_darwin_386
