#!/usr/bin/env bash

set -e

# Determine current operating system.
# If OS is Windows, localhost will be set to VM linux address.
# This only applies to Docker containers.
localhost=127.0.0.1
if [[ "$OSTYPE" == 'msys' ]]; then
    localhost=192.168.99.100
fi

echo "Default localhost address for Docker containers is set to" $localhost

# Export environment variables.
export ADDR=localhost:3000

export SESSIONKEY="secret"

export REDISADDR=$localhost:6379
export DBADDR=$localhost:27017

export REDIS_CONTAINER=redis-server
export MONGO_CONTAINER=mongo-server

if [ "$(docker ps -aq --filter name=$REDIS_CONTAINER)" ]; then
    docker rm -f $REDIS_CONTAINER
fi

# Run Redis Docker container.
docker run \
-d \
--name $REDIS_CONTAINER \
-p 6379:6379 \
redis

if [ "$(docker ps -aq --filter name=$MONGO_CONTAINER)" ]; then
    docker rm -f $MONGO_CONTAINER
fi

# Run Mongo Docker container.
docker run \
-d \
--name $MONGO_CONTAINER \
-p 27017:27017 \
-e MONGO_INITDB_DATABASE=$DBNAME \
drstearns/mongo1kusers

# Run Game API.
go run main.go