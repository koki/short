#!/bin/bash

set -ex

cd $(dirname $0)

./generate_docs.sh
./build.sh
