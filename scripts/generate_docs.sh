#!/bin/bash

set -ex

cd $(dirname $0)

for f in ../docs/*.md
do
	./convert_md_to_text.sh $f
	intermediate=$(echo $f | sed -e "s/.md/.raw.txt/")
	./clean_text_rep.sh $intermediate
done

mkdir -p ../generated

for f in ../docs/*.txt
do
	mv $f ../generated/
done

go-bindata -pkg="generated" -ignore=generated_docs.go$ -o ../generated/generated_docs.go ../generated/
