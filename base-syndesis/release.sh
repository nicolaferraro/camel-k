#!/bin/sh

IMAGE_NAME=quay.io/redhatdemo/camel-k-syndesis-base:latest

docker build -t $IMAGE_NAME .
docker push $IMAGE_NAME
