#!/bin/bash

args=()
for arg in "$@"; do
  args+=("$arg")
done

go run registry/*.go "${args[@]}" > terraform-registry-manifest.json