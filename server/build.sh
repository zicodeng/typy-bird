#!/usr/bin/env bash

set -e

export GATEWAY_CONTAINER=typy-bird-api

# Build the API server Linux executable.
GOOS=linux go build

# Build the API server Docker container image.
docker build -t zicodeng/$GATEWAY_CONTAINER .

if [ "$(docker ps -aq --filter name=$GATEWAY_CONTAINER)" ]; then
    docker rm -f $GATEWAY_CONTAINER
fi

# Remove dangling images.
if [ "$(docker images -q -f dangling=true)" ]; then
    docker rmi $(docker images -q -f dangling=true)
fi

go clean