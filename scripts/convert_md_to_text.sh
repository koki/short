#!/bin/bash

set -ex

input=$1
input_xml=$(echo $input | sed -e "s/.md/.xml/")

mmark -xml2 -page $input > $input_xml 
xml2rfc --raw  $input_xml 
