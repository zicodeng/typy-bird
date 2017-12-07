#!/usr/bin/env bash

export CLIENT_CONTAINER=typy-bird-client

docker pull zicodeng/$CLIENT_CONTAINER

if [ "$(docker ps -aq --filter name=$CLIENT_CONTAINER)" ]; then
    docker rm -f $CLIENT_CONTAINER
fi

docker image prune -f

docker run -d \
-p 80:80 \
--name $CLIENT_CONTAINER \
zicodeng/$CLIENT_CONTAINER
