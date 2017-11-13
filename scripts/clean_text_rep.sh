#!/bin/bash

set -ex

input=$1

sed -i '/Status of This Memo/,+31d' $input
sed -i '/Network Working Group/,+5d' $input 
sed -i '1,3d' $input
sed -i 's/Expires/       /g' $input 
sed -i '/Internet-Draft/,+1d' $input

output=$(echo $input | sed -e "s/.raw.txt/.txt/")
intermediate=$(echo $input | sed -e "s/.raw.txt/.xml/")

mv $input $output
rm -r $intermediate
