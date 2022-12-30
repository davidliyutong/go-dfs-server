package server

import (
	"errors"
	"go-dfs-server/pkg/utils"
	"sync"
	"time"
)

type SessionManager interface {
	New(path string, filePath string, mode int) (string, error)
	Delete(id string) error
	Clean() error
	ListSessions() []Session
	Get(id string) (Session, error)
	SetTimeOut(time.Duration) error
	Reset() error
}

type sessionManager struct {
	sessions sync.Map
	Timeout  time.Duration
}

func (s *sessionManager) Reset() error {
	s.sessions = sync.Map{}
	return nil
}

func (s *sessionManager) New(path string, filePath string, mode int) (string, error) {
	id := utils.MustGenerateUUID()
	s.sessions.Store(id, NewSession(path, filePath, id, mode))
	return id, nil
}

func (s *sessionManager) Delete(id string) error {
	s.sessions.Delete(id)
	return nil
}

func (s *sessionManager) ListSessions() []Session {
	var sessions []Session
	s.sessions.Range(func(k, v interface{}) bool {
		sessions = append(sessions, v.(Session))
		return true
	})
	return sessions
}

func (s *sessionManager) Get(id string) (Session, error) {
	session, ok := s.sessions.Load(id)
	if ok {
		return session.(Session), nil
	} else {
		return nil, errors.New("session not found")
	}
}

func (s *sessionManager) SetTimeOut(duration time.Duration) error {
	s.Timeout = duration
	return nil
}

func (s *sessionManager) Close(id string) error {
	session, ok := s.sessions.Load(id)
	if ok {
		return session.(Session).Close()
	} else {
		return errors.New("session not found")
	}
}

func (s *sessionManager) Clean() error {
	var deadSessions []string
	currTime := time.Now()
	s.sessions.Range(func(k, v interface{}) bool {
		if currTime.Sub(v.(Session).GetTime()) > s.Timeout {
			deadSessions = append(deadSessions, k.(string))
		}
		return true
	})
	for _, id := range deadSessions {
		s.sessions.Delete(id)
	}
	return nil
}

func NewSessionManager() SessionManager {
	return &sessionManager{}
}
