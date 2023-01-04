package main

import (
	"go-dfs-server/pkg/nameserver/server"
	"time"
)

func main() {
	k := server.NewHealthKeeper()
	k.Start()
	time.Sleep(time.Second * 5)
	k.Stop()
	time.Sleep(time.Second * 10)

}
