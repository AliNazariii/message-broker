#!/usr/bin/env bash

./build.sh release
./docker.sh

docker-compose -f deploy/postgres-docker-compose.yml up -d
kubectl create -f deploy/kubernetes.yml