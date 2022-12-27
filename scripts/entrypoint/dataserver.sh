#!/bin/bash
CONFIG_DIR=/config
if [ ! -f $CONFIG_DIR/config.yaml ]; then
    echo "creating config.yaml"
    ./go-dfs-dataserver init --print > "$CONFIG_DIR/config.yaml"
    cat "$CONFIG_DIR/config.yaml"
fi

./go-dfs-dataserver serve --config "$CONFIG_DIR/config.yaml"
