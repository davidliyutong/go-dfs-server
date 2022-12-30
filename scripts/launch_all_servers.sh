#!/bin/bash

DATA_ROOT=$1
TAG=$2
N_SERVERS=$3
NETWORK=$4

if [ -z "$DATA_ROOT" ]; then
    echo "Usage: $0 DATA_ROOT TAG N_SERVERS NETWORK"
    exit 1
fi

if [ -z "$TAG" ]; then
    echo "Usage: $0 DATA_ROOT TAG N_SERVERS NETWORK"
    exit 1
fi

if [ -z "$N_SERVERS" ]; then
    echo "Usage: $0 DATA_ROOT TAG N_SERVERS NETWORK"
    exit 1
fi

if [ -z "$NETWORK" ]; then
    echo "Usage: $0 DATA_ROOT TAG N_SERVERS NETWORK"
    exit 1
fi

if [ ! -d "$DATA_ROOT" ]; then
    echo "DATA_ROOT $DATA_ROOT must be a existing directory"
    exit 1
fi


if [ -z "$(docker network ls | awk '{print $2}' | grep "$NETWORK")" ]; then
    echo "Creating network $NETWORK"
    docker network create "$NETWORK"
fi

DFSAPP_DATA_SERVERS=""

echo "launching $N_SERVERS servers with tag $TAG"
for i in $(seq 0 "$(echo "$N_SERVERS - 1" | bc)"); do
    if [ ! -d "$DATA_ROOT/$i" ]; then
        mkdir -p "$DATA_ROOT/$i"
    fi
    if [ ! -d "$DATA_ROOT/$i/config" ]; then
        mkdir -p "$DATA_ROOT/$i/config"
    fi
    if [ ! -d "$DATA_ROOT/$i/data" ]; then
        mkdir -p "$DATA_ROOT/$i/data"
    fi
    docker stop "DataServer-$i" > /dev/null 2>&1 && docker rm "DataServer-$i" > /dev/null 2>&1
    docker run --rm -d \
           --name "DataServer-$i" \
           -v "$DATA_ROOT"/"$i"/config:/config \
           -v "$DATA_ROOT"/"$i"/data:/data \
           -p "$(echo "27904 + $i" | bc)":27904 \
           --net="$NETWORK" \
           -e "DFSAPP_DEBUG=1" \
           davidliyutong/go-dfs-dataserver:"$TAG"
    DFSAPP_DATA_SERVERS="$DFSAPP_DATA_SERVERS,DataServer-$i:$(echo "27904 + $i" | bc)"
done

if [ ! -d "$DATA_ROOT"/name ]; then
    mkdir -p "$DATA_ROOT"/name/config
    mkdir -p "$DATA_ROOT"/name/data
fi

echo "DFSAPP_DATA_SERVERS=$DFSAPP_DATA_SERVERS"
docker stop "NameServer" > /dev/null 2>&1 && docker rm "NameServer" > /dev/null 2>&1
docker run --rm -d \
       --name "NameServer" \
       -v "$DATA_ROOT"/name/config:/config \
       -v "$DATA_ROOT"/name/data:/data \
       -p 27903:27903 \
       -e "DFSAPP_DATA_SERVERS=$DFSAPP_DATA_SERVERS" \
       --net="$NETWORK" \
       davidliyutong/go-dfs-nameserver:"$TAG"