package main

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"go-dfs-server/pkg/nameserver/server"
	"go-dfs-server/pkg/utils"
	"os"
	"time"
)

func createSessions(m server.SessionManager) {
	s1 := server.NewSession("demo1", "/volume/demo1", utils.MustGenerateUUID(), os.O_RDONLY)
	s2 := server.NewSession("demo2", "/volume/demo2", utils.MustGenerateUUID(), os.O_RDWR)
	_ = s2.SetErrorClose(nil)
	s3 := server.NewSession("demo3", "/volume/demo3", utils.MustGenerateUUID(), os.O_RDONLY)
	_ = s3.SetErrorClose(errors.New("demo"))
	s4 := server.NewSession("demo4", "/volume/demo4", utils.MustGenerateUUID(), os.O_RDONLY)
	_ = s4.SetError(errors.New("demo"))
	_ = m.Add(s1)
	time.Sleep(time.Millisecond * 100)
	_ = m.Add(s2)
	time.Sleep(time.Millisecond * 100)
	_ = m.Add(s3)
	time.Sleep(time.Millisecond * 100)
	_ = m.Add(s4)
	time.Sleep(time.Millisecond * 100)
	time.Sleep(5)
}

func clean(m server.SessionManager) {
	log.Info("Starting session cleaner")
	err := m.SetTimeOut(time.Second * 1) // TODO: make this configurable
	if err != nil {
		return
	}
	go func() {
		for {
			err := m.Clean()
			if err != nil {
				log.Warningln(err)
			}
			log.Debugln("trigger clean, active sessions: ", m.ListSessions())

			time.Sleep(time.Second) // TODO: make this configurable
		}
	}()
}
func main() {
	log.SetLevel(log.DebugLevel)
	m := server.NewSessionManager()
	m.HealthKeeper().Start()
	go clean(m)
	go createSessions(m)

	time.Sleep(time.Second * 200)
	m.HealthKeeper().Stop()
	m.HealthKeeper().Wait()
	fmt.Println("done")

}
