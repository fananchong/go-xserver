#!/bin/bash

set -ex

for plugin_name in `ls -l | grep '^d' | awk '{print $9}'`; do 
    go build $FLAG_RACE -buildmode=plugin -o $BIN_DIR/$plugin_name.so ./$plugin_name;
done

