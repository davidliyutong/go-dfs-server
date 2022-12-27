package main

import (
	"fmt"
	v1 "go-dfs-server/pkg/dataserver/client/v1"
)

func Main() {
	println("Hello, world!")
	cli := v1.NewDataServerClient("", "localhost", 27904, false)
	fmt.Println(cli.Ping())
}
