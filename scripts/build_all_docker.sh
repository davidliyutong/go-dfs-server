#!/bin/bash
TAG=$1

if [ -z "$TAG" ]; then
    TAG=$(git describe --tags)
fi

#docker run --rm -it davidliyutong/go-dfs-dataserver:latest .
docker build -t davidliyutong/go-dfs-dataserver:$TAG --file ./docker/dataserver/Dockerfile .
docker build -t davidliyutong/go-dfs-nameserver:$TAG --file ./docker/nameserver/Dockerfile .
