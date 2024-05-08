#!/bin/sh

docker buildx build --push --platform linux/arm64/v8,linux/amd64 \
  --tag yourrepo/cockroachdb:21.1.9 \
  --build-arg COCKROACHDB_VERSION=21.1.9 .