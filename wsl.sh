#!/bin/bash

set -e

CUR_DIR=$PWD
SRC_DIR=$PWD
BIN_DIR=$PWD/bin
CONF_DIR=$PWD/config

case $1 in
    "start")
        mkdir -p $CONF_DIR
        ln -sf $SRC_DIR/common/config/framework.toml $CONF_DIR/
        ln -sf $SRC_DIR/default_plugins/login/login.toml $CONF_DIR/
        cd $BIN_DIR
        mkdir -p $BIN_DIR/logs
        mkdir -p $BIN_DIR/logs.back
        if [[ `ls -l ./logs | wc -l` > 1  ]]; then
            mv -f ./logs/* ./logs.back/
        fi
        nohup ./go-xserver --app mgr --network-port '0,31000' --common-logflushinterval 200 > /dev/null 2>&1 &
        nohup ./go-xserver --app login --network-port '7200,0' --common-logflushinterval 200 --suffix 1 > /dev/null 2>&1 &
        nohup ./go-xserver --app login --network-port '7201,0' --common-logflushinterval 200 --suffix 2 > /dev/null 2>&1 &
        nohup ./go-xserver --app login --network-port '7202,0' --common-logflushinterval 200 --suffix 3 > /dev/null 2>&1 &
        nohup ./go-xserver --app gateway --network-port '7300,33000' --common-logflushinterval 200 --suffix 1 > /dev/null 2>&1 &
        nohup ./go-xserver --app gateway --network-port '7301,33001' --common-logflushinterval 200 --suffix 2 > /dev/null 2>&1 &
        nohup ./go-xserver --app gateway --network-port '7302,33002' --common-logflushinterval 200 --suffix 3 > /dev/null 2>&1 &
        nohup ./go-xserver --app lobby --network-port '0,0' --common-logflushinterval 200 --suffix 1 > /dev/null 2>&1 &
        nohup ./go-xserver --app lobby --network-port '0,0' --common-logflushinterval 200 --suffix 2 > /dev/null 2>&1 &
        nohup ./go-xserver --app lobby --network-port '0,0' --common-logflushinterval 200 --suffix 3 > /dev/null 2>&1 &
        nohup ./go-xserver --app match --network-port '0,0' --common-logflushinterval 200 --suffix 1 > /dev/null 2>&1 &
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
esac

echo "Usage:"
echo "    wsl.sh start"
echo "    wsl.sh stop"

cd $CUR_DIR

exit 0

