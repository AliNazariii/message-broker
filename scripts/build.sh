#!/usr/bin/env bash

BROKER_VERSION=`${PWD}/scripts/version.sh`
TIME=$(date)

if [[ -d build/release ]]; then
    mkdir -p build/release
    mkdir -p build/debug
fi

if [[ "$1" == "release" ]]; then
    echo "Building in release mode"
#    go build -o build/release/nazari-broker_${BROKER_VERSION} -a -installsuffix cgo -ldflags="-X 'broker/version.BuildTime=$TIME' -X 'broker/version.BuildVersion=$BROKER_VERSION' -s" main.go
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/release/nazari-broker_"${BROKER_VERSION}" -a -installsuffix cgo -ldflags="-X 'broker/version.BuildTime=$TIME' -X 'broker/version.BuildVersion=$BROKER_VERSION' -s" main.go
else
    echo "Building in debug mode"
    # shellcheck disable=SC2086
    go build -o build/debug/nazari-broker_${BROKER_VERSION} -a -v -installsuffix cgo -ldflags="-X 'broker/version.BuildTime=$TIME' -X 'broker/version.BuildVersion=$BROKER_VERSION' -s" main.go
fi