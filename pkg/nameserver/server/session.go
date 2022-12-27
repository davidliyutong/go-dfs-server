package server

import (
	"errors"
	v1 "go-dfs-server/pkg/nameserver/apiserver/blob/v1/model"
	"go-dfs-server/pkg/utils"
	"sync"
	"time"
)

type SessionManager interface {
	New(path string) (string, error)
	Delete(id string) error
	Clean() error
	ListSessions() []v1.Session
	Get(id string) (v1.Session, error)
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

func (s *sessionManager) New(path string) (string, error) {
	id := utils.MustGenerateUUID()
	s.sessions.Store(id, v1.NewSession(path, id))
	return id, nil
}

func (s *sessionManager) Delete(id string) error {
	s.sessions.Delete(id)
	return nil
}

func (s *sessionManager) ListSessions() []v1.Session {
	var sessions []v1.Session
	s.sessions.Range(func(k, v interface{}) bool {
		sessions = append(sessions, v.(v1.Session))
		return true
	})
	return sessions
}

func (s *sessionManager) Get(id string) (v1.Session, error) {
	session, ok := s.sessions.Load(id)
	if ok {
		return session.(v1.Session), nil
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
		return session.(v1.Session).Close()
	} else {
		return errors.New("session not found")
	}
}

func (s *sessionManager) Clean() error {
	var deadSessions []string
	currTime := time.Now()
	s.sessions.Range(func(k, v interface{}) bool {
		if currTime.Sub(v.(v1.Session).GetTime()) > s.Timeout {
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
