#!/bin/bash

set -e

cd $(dirname $0)/..

mkdocs build --clean

aws s3 sync site/ s3://docs.koki.io/short/ --delete --acl public-read
