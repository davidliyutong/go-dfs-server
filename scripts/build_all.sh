#!/bin/bash

if [[ ! -d "build" ]]; then
  mkdir build
fi

go build -o ./build/nameserver ./cmd/nameserver/app.go
go build -o ./build/dataserver ./cmd/dataserver/app.go
go build -o ./build/client ./cmd/client/app.go

