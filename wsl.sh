#!/bin/bash

set -e

CUR_DIR=$PWD
BIN_DIR=$PWD/bin

case $1 in
    "start")
        cd $BIN_DIR
        mkdir -p $BIN_DIR/logs
        mkdir -p $BIN_DIR/logs.back
        if [[ `ls -l ./logs | wc -l` > 1  ]]; then
            mv -f ./logs/* ./logs.back/
        fi
        nohup ./go-xserver --app mgr --network-port '0,30000' > /dev/null 2>&1 &
        nohup ./go-xserver --app login --network-port '7500,0' --suffix 1 > /dev/null 2>&1 &
        nohup ./go-xserver --app login --network-port '7501,0' --suffix 2 > /dev/null 2>&1 &
        nohup ./go-xserver --app login --network-port '7502,0' --suffix 3 > /dev/null 2>&1 &
        nohup ./go-xserver --app gateway --network-port '7600,36001' --suffix 1 > /dev/null 2>&1 &
        nohup ./go-xserver --app gateway --network-port '7601,36002' --suffix 2 > /dev/null 2>&1 &
        nohup ./go-xserver --app gateway --network-port '7602,36003' --suffix 3 > /dev/null 2>&1 &
        nohup ./go-xserver --app lobby --network-port '7700,0' --suffix 1 > /dev/null 2>&1 &
        nohup ./go-xserver --app lobby --network-port '7701,0' --suffix 2 > /dev/null 2>&1 &
        nohup ./go-xserver --app lobby --network-port '7702,0' --suffix 3 > /dev/null 2>&1 &
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

