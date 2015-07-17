#!/bin/bash

set -e
set -u

for line in "$(cat Godeps)"; do
    importpath=$(echo "$line" | awk '{print $1}')
    revision=$(echo "$line" | awk '{print $2}')
    gb vendor fetch -revision "$revision" "$importpath"
done
