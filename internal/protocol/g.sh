#!/bin/bash

set -ex

docker run --rm -v $PWD:/out -w /out znly/protoc --gogofaster_out=. -I=. *.proto
