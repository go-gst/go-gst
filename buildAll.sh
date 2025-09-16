#!/usr/bin/env bash

# this script builds all go packages in the current directory

packages=$(go list ./...)
for package in $packages; do
    go build -o /dev/null "$package" || exit 1
done

packages=$(go list ./pkg/...)
for package in $packages; do
    go build -o /dev/null "$package" || exit 1
done
