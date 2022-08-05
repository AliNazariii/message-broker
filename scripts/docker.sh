#!/usr/bin/env bash

BROKER_VERSION=`${PWD}/scripts/version.sh`
TIME=$(date)

echo "version: ${BROKER_VERSION}"
docker build --build-arg docker_version="${BROKER_VERSION}" -t docker.bale.ai/bale/nazari-broker:"${BROKER_VERSION}" -f deploy/Dockerfile .
docker push docker.bale.ai/bale/nazari-broker:"${BROKER_VERSION}"