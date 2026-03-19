#!/usr/bin/env bash

# This script downloads the gir-files-rs repository and copies the GIR files into the
# girs directory

set -euo pipefail

repository="https://gitlab.freedesktop.org/gstreamer/gir-files-rs.git"

destination="girs"

if [ ! -d "$destination" ]; then
    echo "Creating directory $destination"
    mkdir "$destination"
fi

# clone into a temporary directory
temp_dir=$(mktemp -d)
trap 'rm -rf "$temp_dir"' EXIT


git clone "$repository" "$temp_dir"


cp -r "$temp_dir"/*.gir "$destination"
