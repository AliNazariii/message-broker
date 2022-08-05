#!/usr/bin/env bash
#
# Generate all chie protobuf bindings.
# Run from repository root.
#
set -e

if ! [[ "$0" =~ scripts/genproto.sh ]]; then
	echo "must be run from repository root"
	exit 255
fi

IFS=. V=($(protoc --version | cut -f2 -d' '))

if [[ ${V[0]} -lt 3 || (${V[0]} == 3 && ${V[1]} -lt 7) ]]; then
	echo "could not find protoc version grater than 3.7.1, is it installed + in PATH?"
	exit 255
fi

cd ${PWD}/api/proto
rm -rf src
mkdir src
protoc --go_out=plugins=grpc:src ./*.proto
cd ..
