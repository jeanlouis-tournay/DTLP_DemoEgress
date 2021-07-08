#!/usr/bin/env bash
# Docker registry where image will be pushed
if [ -z "$DOCKER_REGISTRY_URL" ]; then
    echo "Docker Registry URL not set"
    exit 1
fi
# Docker library of image repository
#if [ -z "$DOCKER_PUSH_LIBRARY" ]; then
#    echo "Docker library not set"
#    exit 1
#fi

# Service to make
if [ -z "$1" ]; then
    echo "Service not set"
    exit 1
fi
export SERVICE_NAME=$1
make $2