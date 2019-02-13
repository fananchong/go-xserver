#!/bin/bash

set -ex

docker run --rm -v $PWD:/out -w /out znly/protoc --gogofaster_out=. -I=. *.proto
docker run --rm -v $PWD:/out -w /out znly/protoc --python_out=. -I=. *.proto
sed -i 's/import lobby_custom_pb2/import protocol.lobby_custom_pb2/g' ./lobby_pb2.py

