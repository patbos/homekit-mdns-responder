#!/bin/bash

echo "Environment:"
echo "----------------------------"
env
echo "----------------------------"

if [ "$TRAVIS_PULL_REQUEST" = "true" ] || [ "$TRAVIS_BRANCH" != "master" ]; then
  docker buildx build \
    --progress plain \
    --platform=linux/amd64,linux/arm64,linux/arm/v7,linux/arm/v6 \
    .
  exit $?
fi
echo $DOCKER_PASSWORD | docker login -u $DOCKER_USER --password-stdin &> /dev/null
TAG="${TRAVIS_TAG:-latest}"
docker buildx build \
     --progress plain \
    --platform=linux/amd64,linux/arm64,linux/arm/v7,linux/arm/v6 \
    -t $DOCKER_REPO:$TAG \
    --push \
    .
