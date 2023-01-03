package main

import (
	log "github.com/sirupsen/logrus"
	"go-dfs-server/pkg/nameserver/server"
)

func main() {
	manager := server.NewSessionManager()
	s, err := manager.New("./test", "", 0)
	if err != nil {
		return
	}
	log.Info(s)
	session, _ := manager.Get("./test")
	log.Info(session)
}
