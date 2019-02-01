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

cd $SRC_DIR
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
    "start")
        cd $BIN_DIR
        mkdir -p $BIN_DIR/logs
        mkdir -p $BIN_DIR/logs.back
        if [[ `ls -l ./logs | wc -l` > 1  ]]; then
            mv -f ./logs/* ./logs.back/
        fi
        nohup ./go-xserver --app mgr --network-port '0,30000' > /dev/null 2>&1 &
        nohup ./go-xserver --app login --network-port '7500,0' > /dev/null 2>&1 &
        nohup ./go-xserver --app gateway --network-port '7501,30001' > /dev/null 2>&1 &
        nohup ./go-xserver --app lobby --network-port '7502,0' > /dev/null 2>&1 &
        ps -ux | grep go-xserver
        ;;
    ?);;
esac

cd $CUR_DIR

exit 0

