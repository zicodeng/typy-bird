#!/usr/bin/env bash

set -e

export GATEWAY_CONTAINER=typy-bird-api
export MONGO_CONTAINER=mongo-server

export APP_NETWORK=appnet

export ADDR=:80
export DBADDR=$MONGO_CONTAINER:27017

# Make sure to get the latest image.
docker pull zicodeng/$GATEWAY_CONTAINER

# Remove the old containers first.
if [ "$(docker ps -aq --filter name=$GATEWAY_CONTAINER)" ]; then
    docker rm -f $GATEWAY_CONTAINER
fi

if [ "$(docker ps -aq --filter name=$MONGO_CONTAINER)" ]; then
    docker rm -f $MONGO_CONTAINER
fi

# Remove dangling images.
if [ "$(docker images -q -f dangling=true)" ]; then
    docker rmi $(docker images -q -f dangling=true)
fi

# Clean up the system.
docker system prune -f

# Create Docker private network if not exist.
if ! [ "$(docker network ls | grep $APP_NETWORK)" ]; then
    docker network create $APP_NETWORK
fi

# Run Mongo Docker container inside our appnet private network.
docker run \
-d \
--name mongo-server \
--network $APP_NETWORK \
--restart unless-stopped \
mongo

# Run Info 344 API Gateway Docker container inside our appnet private network.
docker run \
-d \
-p 80:80 \
--name $GATEWAY_CONTAINER \
--network $APP_NETWORK \
-e ADDR=$ADDR \
-e DBADDR=$DBADDR \
--restart unless-stopped \
zicodeng/$GATEWAY_CONTAINER