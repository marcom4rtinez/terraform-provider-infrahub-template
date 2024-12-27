#!/bin/bash

args=()
for arg in "$@"; do
  args+=("$arg")
done

search_string="\"version\": \"$2\","

if grep -q "$search_string" terraform-registry-manifest.json; then
echo "already up2date! Adding hashes now..."
cat dist/*SHA256SUMS | go run registry/*.go --hashes --manifest terraform-registry-manifest.json
else
echo "wrong version found! Regenerating..."
go run registry/*.go "${args[@]}" > terraform-registry-manifest.json
fi