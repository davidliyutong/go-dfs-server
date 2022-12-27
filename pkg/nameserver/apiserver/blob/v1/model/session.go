package v1

import (
	"sync"
	"time"
)

type Session interface {
	Open(mode int) error
	Seek(offset int64, whence int) (int64, error)
	Close() error
	Flush() error
	Read(size int64) ([]byte, error)
	Write(data []byte) error
	Truncate(size int64) error
	Lock() error
	Unlock() error
	GetTime() time.Time
}

type session struct {
	Path    string
	ID      string
	Time    time.Time
	Offset  int64
	IsAlive bool
	Blob    BlobStruct
	Mutex   sync.RWMutex
}

func (s *session) Truncate(size int64) error {
	//TODO implement me
	panic("implement me")
}

func (s *session) Write(data []byte) error {
	//TODO implement me
	panic("implement me")
}

func (s *session) Open(mode int) error {
	//TODO implement me
	panic("implement me")
}

func (s *session) Seek(offset int64, whence int) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (s *session) Close() error {
	//TODO implement me
	panic("implement me")
}

func (s *session) Flush() error {
	//TODO implement me
	panic("implement me")
}

func (s *session) Read(size int64) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (s *session) Lock() error {
	//TODO implement me
	panic("implement me")
}

func (s *session) Unlock() error {
	//TODO implement me
	panic("implement me")
}

func (s *session) GetTime() time.Time {
	return s.Time
}

func NewSession(path string, id string) Session {
	return &session{
		Path:    path,
		ID:      id,
		Time:    time.Now(),
		Offset:  0,
		IsAlive: true,
	}
}
