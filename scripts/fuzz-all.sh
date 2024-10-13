#!/bin/bash

set -e

fuzzTime=${1:-10}

files=$(grep -r --include='**_test.go' --files-with-matches 'func Fuzz' .)

echo "Fuzz time: ${fuzzTime}s"
echo "Fuzz targets: ${files}"

for file in ${files}; do
	funcs=$(grep -oP 'func \K(Fuzz\w*)' $file)
	for func in ${funcs}; do
		echo "Fuzzing $func in $file"
		parentDir=$(dirname $file)
		go test -tags='fuzz' $parentDir -run=Fuzz -fuzz=$func -fuzztime=${fuzzTime}s
	done
done
