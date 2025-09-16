#!/usr/bin/env bash

# this script builds all go packages in the current directory

packages=$(go list ./pkg/...)
for package in $packages; do
    echo "building $package" 
    go build -o /dev/null "$package" || exit 1
done

examples=$(go list ./examples/...)
for ex in $examples; do
    echo "building example $ex" 
    go build -o /dev/null "$ex" || exit 1
done
