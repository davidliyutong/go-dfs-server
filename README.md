# Distributed File Server written in GO


# How to build

To build docker images, run

```shell
make image
```

This will build the dataserver and nameserver with the latest git tag.

# Get Started

To run the demo with docker, run

```shell
make demo
```

> This command pulls image from DockerHub, thus need active Internet connection

To build image from source, run

```shell
make demo.prepare
make demo.start
```

To stop demo, run

```shell
make demo.stop
```