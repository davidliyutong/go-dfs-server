package server

import (
	"errors"
	"go-dfs-server/pkg/utils"
	"sync"
	"time"
)

type SessionManager interface {
	New(path string, filePath string, mode int) (Session, error)
	Delete(path string) error
	Clean() error
	ListSessions() []Session
	Get(path string) (Session, error)
	SetTimeOut(time.Duration) error
	Reset() error
}

type sessionManager struct {
	sessions sync.Map
	timeout  time.Duration
}

func (s *sessionManager) Reset() error {
	s.sessions = sync.Map{}
	s.timeout = time.Minute * 30
	return nil
}

func (s *sessionManager) New(path string, filePath string, mode int) (Session, error) {
	_, ok := s.sessions.Load(path)
	if ok {
		return nil, errors.New("file already opened")
	} else {
		session := NewSession(path, filePath, utils.MustGenerateUUID(), mode)
		s.sessions.Store(path, session)
		return session, nil
	}

}

func (s *sessionManager) Delete(path string) error {
	session, ok := s.sessions.Load(path)
	if ok {
		session.(Session).Delete()
		session.(Session).Wait()
		s.sessions.Delete(path)
		return nil
	} else {
		return errors.New("session not found")
	}
}

func (s *sessionManager) ListSessions() []Session {
	var sessions []Session
	s.sessions.Range(func(k, v interface{}) bool {
		sessions = append(sessions, v.(Session))
		return true
	})
	return sessions
}

func (s *sessionManager) Get(path string) (Session, error) {
	session, ok := s.sessions.Load(path)
	if ok {
		return session.(Session), nil
	} else {
		return nil, errors.New("session not found")
	}
}

func (s *sessionManager) SetTimeOut(duration time.Duration) error {
	s.timeout = duration
	return nil
}

func (s *sessionManager) Close(path string) error {
	session, ok := s.sessions.Load(path)
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
		if currTime.Sub(v.(Session).GetTime()) > s.timeout {
			deadSessions = append(deadSessions, k.(string))
		}
		return true
	})
	for _, path := range deadSessions {
		path := path

		err := s.Delete(path)
		if err != nil {
			// TODO: handle error
		}
	}
	return nil
}

func NewSessionManager() SessionManager {
	return &sessionManager{
		sessions: sync.Map{},
		timeout:  time.Minute * 30,
	}
}
