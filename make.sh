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
        for plugin_name in $plugins; do
            c=1
            if [[ $plugin_name != "mgr" ]]; then
                c=3
            fi
            for (( i=1; i<=$c; i++ )); do
                nohup ./go-xserver --app $plugin_name > /dev/null 2>&1 &
            done
        done
        ps -ux | grep go-xserver
        ;;
    ?);;
esac

cd $CUR_DIR

exit 0

