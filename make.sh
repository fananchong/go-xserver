#!/bin/bash

set -ex


if [[ $1 == "stop" ]]; then
    pkill go-xserver
    ps -ux | grep go-xserver
    exit 0 
fi


CUR_DIR=$PWD
SRC_DIR=$PWD
BIN_DIR=$PWD/bin
SERVICE_DIR=$PWD/services/
FLAGS=-race
export GOBIN=$BIN_DIR

go generate ./...

cd $SRC_DIR
go vet ./...

cd $SERVICE_DIR
plugins=`ls -l | grep '^d' | awk '{print $9}' | grep -v 'internal'`
for plugin_name in $plugins; do
    go build $FLAGS -buildmode=plugin -o $BIN_DIR/$plugin_name.so ./$plugin_name;
done
cd $SRC_DIR
go install $FLAGS .

case $1 in
    "docker") docker build -t go-xserver . ;;
    "start")
        cd $BIN_DIR
        mkdir -p $BIN_DIR/logs
        mkdir -p $BIN_DIR/logs.back
        mv $BIN_DIR/logs/* $BIN_DIR/logs.back/
        for plugin_name in $plugins; do
            nohup ./go-xserver --app $plugin_name > /dev/null 2>&1 &
        done
        ps -ux | grep go-xserver
        ;;
    ?);;
esac

cd $CUR_DIR

exit 0

