package repo

import (
	"go-dfs-server/pkg/nameserver/server"
)

type BlobRepo interface {
	Open(path string, mode int) (string, error) // Sync
	Close(sessionID string) error               // Sync
	Flush(sessionID string) error               // Sync
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
