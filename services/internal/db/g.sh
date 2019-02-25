#!/bin/bash

set -ex

docker run --rm -v $PWD/redis_def:/app/input -v $PWD:/app/output fananchong/redis2go --package=db

