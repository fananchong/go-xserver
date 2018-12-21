#!/bin/bash

set -ex

SRC_DIR=$PWD
BIN_DIR=$PWD/bin
SERVICE_DIR=$PWD/services/
PROTO_DIR=$PWD/internal/protocol/
FLAGS=-race
GOBIN=$BIN_DIR

cd $PROTO_DIR
./g.sh

cd $SRC_DIR
go vet ./...

cd $SERVICE_DIR
for plugin_name in `ls -l | grep '^d' | awk '{print $9}'`; do
    go build $FLAGS -buildmode=plugin -o $BIN_DIR/$plugin_name.so ./$plugin_name;
done
cd $SRC_DIR
go install $FLAGS .

case $1 in
    "docker") docker build -t go-xserver . ;;
    "start")
        cd $BIN_DIR
        ;;
    ?);;
esac
