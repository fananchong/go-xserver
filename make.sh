#!/bin/bash

set -e

CUR_DIR=$PWD
SRC_DIR=$PWD
BIN_DIR=$PWD/bin
CONF_DIR=$PWD/config
SERVICE_DIR=$PWD/services/
FLAGS=-race

case $1 in
    "start")
        cd $SERVICE_DIR
        plugins=`ls -l | grep '^d' | awk '{print $9}' | grep -v 'internal'`
        mkdir -p $CONF_DIR
        ln -sf $SRC_DIR/common/config/framework.toml $CONF_DIR/
        ln -sf $SRC_DIR/default_plugins/login/login.toml $CONF_DIR/
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
                nohup ./go-xserver --app $plugin_name --suffix $i > /dev/null 2>&1 &
            done
        done
        sleep 1s
        ps -ux | grep go-xserver
        exit 0
        ;;
    "stop")
        pkill go-xserver
        sleep 5s
        ps -ux | grep go-xserver
        exit 0
        ;;
    "")
        export GOPROXY=https://goproxy.io
        cd $SRC_DIR
        go generate ./...
        cd $SRC_DIR
        go vet ./...
        echo "start build ..."
        cd $SRC_DIR/default_plugins
        plugins=`ls -l | grep '^d' | awk '{print $9}' | grep -v 'internal'`
        for plugin_name in $plugins; do
            go build $FLAGS -buildmode=plugin -o $BIN_DIR/$plugin_name.so ./$plugin_name;
        done
        cd $SERVICE_DIR
        plugins=`ls -l | grep '^d' | awk '{print $9}' | grep -v 'internal'`
        for plugin_name in $plugins; do
            go build $FLAGS -buildmode=plugin -o $BIN_DIR/$plugin_name.so ./$plugin_name;
        done
        cd $SRC_DIR
        go build $FLAGS -o $BIN_DIR/go-xserver .
        echo "done"
        exit 0
        ;;
esac

echo "Usage:"
echo "    make.sh"
echo "    make.sh start"
echo "    make.sh stop"

cd $CUR_DIR

exit 0

