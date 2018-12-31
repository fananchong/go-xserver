#!/bin/bash

set -ex

CUR_DIR=$PWD
SRC_DIR=$PWD
BIN_DIR=$PWD/bin
SERVICE_DIR=$PWD/services/
PROTO_DIR=$PWD/internal/protocol/
DB_DIR=$PWD/internal/db/
FLAGS=-race
export GOBIN=$BIN_DIR

if [[ $1 == "prebuild" ]]; then
    cd $PROTO_DIR
    ./g.sh
    cd $DB_DIR
    ./g.sh
fi

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

cd $CUR_DIR
