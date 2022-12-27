FROM golang:1.19.1 as builder

WORKDIR /opt/build
COPY . .

ENV GOPROXY https://goproxy.cn


RUN make build

FROM ubuntu:20.04

WORKDIR /opt/app

COPY --from=builder /opt/build/_output/platforms/linux/amd64/go-dfs-dataserver ./go-dfs-dataserver

COPY scripts/entrypoint/dataserver.sh ./entrypoint.sh
RUN chmod +x ./entrypoint.sh

VOLUME /data
VOLUME /config

SHELL ["/bin/bash", "-c"]
ENTRYPOINT ["./entrypoint.sh"]