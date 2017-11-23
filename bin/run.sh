#!/usr/bin/env bash

set -e

TAG="latest"
SLACK_TOKEN=$(echo $SLACK_TOKEN)
NAME="bitcoin-bot"

docker build . -t $TAG
docker run --name $NAME --env SLACK_TOKEN=$SLACK_TOKEN -p 8080:8080 -it $NAME:$TAG