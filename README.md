# Distributed File Server written in GO


# How to build

To build docker images, run the build script

```shell
./scripts/build_all_docker.sh <tag>
```

This will build the dataserver and nameserver with the specified tag. 

> If no tag is provided, the script will use `$(git describe --tags)` as build tag.

# Get Started

To run the demo with docker, use the launch script

```shell
./scripts/launch_all_servers.sh <volume> <tab> <number of dataserver> <docker network name>
```

For example, to use the latest
```shell
mkdir data
./scripts/launch_all_servers.sh $(pwd)/data $(git describe --tags) 4 dfs
```