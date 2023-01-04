package server

import (
	"context"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

type HealthTask struct {
	Path    string
	ChunkID int64
	Msg     string
}

type healthKeeper struct {
	sessions    sync.Map
	sessionChan chan Session
	taskGroup   *sync.WaitGroup

	stopCtx context.Context
	stopFn  func()
}

func (h *healthKeeper) router(s Session) {
	log.Debugln("healKeeper routing session: ", s)

	*s.Error() = nil
	h.Done()
	h.sessions.Delete(*s.Path())
	// TODO: finish this router to handle:
	// 1. DataServerOffline
	// 2. Checksum mismatch

}

func (h *healthKeeper) Start() {
	h.taskGroup.Wait()
	h.stopCtx, h.stopFn = context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-h.stopCtx.Done():
				h.taskGroup.Wait()
				return
			default:
				s, ok := <-h.sessionChan
				if ok {
					log.Debugf("add session %v:%v to health keeper", s.Path(), s.ID())
					h.sessions.Store(s.Path(), s)
					go h.router(s)
				}
			}
		}
	}()
}

func (h *healthKeeper) Stop() {
	close(h.sessionChan)
	time.Sleep(time.Second * 3)
	h.stopFn()
	h.taskGroup.Wait()
}

func (h *healthKeeper) Add(s Session) {
	h.sessionChan <- s
	h.taskGroup.Add(1)
}

func (h *healthKeeper) Done() {
	h.taskGroup.Done()
}

func (h *healthKeeper) Wait() {
	h.taskGroup.Wait()
}

func (h *healthKeeper) GetSessions() []Session {
	var sessions []Session
	h.sessions.Range(func(k, v interface{}) bool {
		sessions = append(sessions, v.(Session))
		return true
	})
	return sessions
}

type HealthKeeper interface {
	// Start starts the health keeper.
	Start()
	// Stop stops the health keeper.
	Stop()
	// Add adds a new health message to the health keeper.
	Add(s Session)
	Done()
	Wait()
	GetSessions() []Session
}

var _ HealthKeeper = (*healthKeeper)(nil)

func NewHealthKeeper() HealthKeeper {
	return &healthKeeper{
		sessions:    sync.Map{},
		sessionChan: make(chan Session, 16),
		taskGroup:   &sync.WaitGroup{},
	}
}
