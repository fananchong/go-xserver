#!/bin/bash

set -ex

export SRC_DIR=$PWD
export BIN_DIR=$PWD/bin
export SERVICE_DIR=$PWD/services/
export FLAG_RACE=-race
export GOBIN=$BIN_DIR

go vet ./...
for dir in $SERVICE_DIR; do
    if [[ $dir ]]; then
        cd $dir && ./build.sh
    fi
done
cd $SRC_DIR
go install $FLAG_RACE .

case $1 in
    "docker") docker build -t go-xserver . ;;
    "start")
        cd $BIN_DIR
        ;;
    ?);;
esac

set +ex
export SRC_DIR=
export BIN_DIR=
export SERVICE_DIR=
export FLAG_RACE=
export GOBIN=
