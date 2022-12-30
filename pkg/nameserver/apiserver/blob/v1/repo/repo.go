package repo

import (
	"go-dfs-server/pkg/nameserver/server"
)

type BlobRepo interface {
	Open(path string, mode int) (string, error)
	Close(sessionID string) error
	Flush(sessionID string) error
	Lock(sessionID string) error
	GetLock(path string) ([]string, error)
	LockUnique(sessionID string) error
	Unlock(sessionID string) error
	SessionManager() server.SessionManager
	DataServerManager() server.DataServerManager
}

type Repo interface {
	BlobRepo() BlobRepo
	Close() error
}

var client Repo

func Client() Repo {
	return client
}

func SetClient(c Repo) {
	client = c
}
