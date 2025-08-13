#!/bin/sh

IMAGE_NAME="lo_task"

docker build -f ./deploy/Dockerfile -t $IMAGE_NAME .

docker run --rm --env-file .env -p 9090:9090 $IMAGE_NAME