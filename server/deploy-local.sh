#!/bin/bash
echo "Deleting game-api container..."
docker rm -f game-api

echo "Deleting redis container..."
docker rm -f redisdev

echo "Deleting mongodb container..."
docker rm -f mongodev

echo "Deleting game-network..."
docker network rm game-network

echo "Starting game-network..."
docker network create game-network

echo "Starting redis docker..."
docker run -d --name redisdev --network game-network redis

echo "Starting mongo docker"
docker run -d --name mongodev --network game-network mongo

echo "Starting game-api..."
docker run -d -p 80:80 --name game-api --network game-network -e SESSIONKEY=development -e REDISADDR=redisdev:6379 -e DBADDR=mongodev:27017 bond00729/game-api

# -v $(pwd)/keys:/keys:ro -e TLSCERT=/keys/fullchain.pem -e TLSKEY=/keys/privkey.pem