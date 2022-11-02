#!/bin/bash

if [[ ! -d "build" ]]; then
  mkdir build
fi

go build -o ./build/go-dfs-nameserver ./cmd/nameserver/app.go
go build -o ./build/go-dfs-dataserver ./cmd/dataserver/app.go
go build -o ./build/go-dfs-client ./cmd/client/app.go

