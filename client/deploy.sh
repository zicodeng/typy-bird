#!/usr/bin/env bash

set -e

./build.sh

export CILENT_CONTAINER=typy-bird-client

docker push zicodeng/$CILENT_CONTAINER

ssh -oStrictHostKeyChecking=no root@198.199.105.119 'bash -s' < run.sh