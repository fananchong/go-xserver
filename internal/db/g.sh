#!/bin/bash

set -ex

docker run --rm -v $PWD/../:/out -w /out/db znly/protoc --gogofaster_out=. -I=. -I=../protocol *.proto
sed -i 's#import protocol "."#import protocol "github.com/fananchong/go-xserver/internal/protocol"#g' ./token.pb.go
docker run --rm -v $PWD/redis_def:/app/input -v $PWD:/app/output fananchong/redis2go --package=db
