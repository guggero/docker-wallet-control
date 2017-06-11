#!/bin/bash

LATEST_TAG=$(git describe --abbrev=0 --tags)
DOCKER_NAME=guggero/docker-wallet-control

docker build -t $DOCKER_NAME -t $DOCKER_NAME:latest -t $DOCKER_NAME:$LATEST_TAG .
docker push $DOCKER_NAME
docker push $DOCKER_NAME:latest
docker push $DOCKER_NAME:$LATEST_TAG
