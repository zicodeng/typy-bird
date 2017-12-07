#!/usr/bin/env bash

set -e

export GATEWAY_CONTAINER=typy-bird-api

./build.sh

docker push zicodeng/$GATEWAY_CONTAINER

ssh -oStrictHostKeyChecking=no root@162.243.129.167 'bash -s' < run.sh