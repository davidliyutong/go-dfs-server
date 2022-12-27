FROM golang:1.19.1 as builder

WORKDIR /opt/build
COPY . .

ENV GOPROXY https://goproxy.cn


RUN go mod vendor && go build -o go-dfs-nameserver ./cmd/nameserver/app.go

FROM ubuntu:20.04

WORKDIR /opt/app

COPY --from=builder /opt/build/go-dfs-nameserver ./go-dfs-nameserver

COPY scripts/entrypoint/nameserver.sh ./entrypoint.sh
RUN chmod +x ./entrypoint.sh

VOLUME /data
VOLUME /config

SHELL ["/bin/bash", "-c"]
ENTRYPOINT ["./entrypoint.sh"]