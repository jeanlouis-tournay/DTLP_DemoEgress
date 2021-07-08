# syntax = docker/dockerfile:1.0-experimental

ARG DOCKER_REGISTRY_URL

FROM ${DOCKER_REGISTRY_URL}/product-deployment-base:latest
ARG SERVICE_NAME
COPY --chown=eurocontrol:eurocontrol dist/artifacts/linux/$SERVICE_NAME/app bin/
RUN chmod +x bin/*
