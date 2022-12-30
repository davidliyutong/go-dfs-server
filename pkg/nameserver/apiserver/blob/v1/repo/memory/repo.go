package memory

import (
	"errors"
	repo2 "go-dfs-server/pkg/nameserver/apiserver/blob/v1/repo"
	"go-dfs-server/pkg/nameserver/server"
	"os"
	"sync"
)

type blobRepo struct {
	servers  server.DataServerManager
	sessions server.SessionManager
	locks    server.LockManager
}

var _ repo2.BlobRepo = &blobRepo{nil, nil, nil}

func newBlobRepo(servers server.DataServerManager, sessions server.SessionManager, locks server.LockManager) repo2.BlobRepo {
	return &blobRepo{
		servers:  servers,
		sessions: sessions,
		locks:    locks,
	}
}

func (r *blobRepo) Close(sessionID string) error {
	session, _ := r.SessionManager().Get(sessionID)
	if session == nil {
		return errors.New("session not found")
	} else {
		err := session.Close()
		if err != nil {
			return err
		} else {
			if session.IsOpened() {
				return errors.New("failed to close session")
			} else {
				_ = r.Unlock(sessionID)
				return r.SessionManager().Delete(sessionID)
			}
		}
	}
}

func (r *blobRepo) Open(path string, mode int) (string, error) {

	sessionID, _ := r.SessionManager().New(path, "", mode)
	session, _ := r.SessionManager().Get(sessionID)

	err := session.Open()
	if err != nil {
		_ = r.Close(sessionID)
		return "", err
	} else {
		switch mode {
		case os.O_RDONLY:
			err = r.Lock(sessionID)
		case os.O_RDWR:
			err = r.LockUnique(sessionID)
		default:
		}
		if err != nil {
			_ = r.Close(sessionID)
			return "", err
		} else {
			return sessionID, nil
		}
	}
}

func (r *blobRepo) Flush(sessionID string) error {
	session, _ := r.SessionManager().Get(sessionID)
	if session == nil {
		return errors.New("session not found")
	} else {
		return session.Flush()
	}
}

func (r *blobRepo) Lock(sessionID string) error {
	session, err := r.SessionManager().Get(sessionID)
	if err != nil {
		return err
	}
	if session.IsOpened() {
		return r.locks.Lock(*session.GetPath(), sessionID)
	} else {
		return errors.New("session not opened")
	}
}

func (r *blobRepo) GetLock(path string) ([]string, error) {
	return r.locks.GetLock(path)
}
func (r *blobRepo) LockUnique(sessionID string) error {
	session, err := r.SessionManager().Get(sessionID)
	if err != nil {
		return err
	}
	if session.IsOpened() {
		return r.locks.LockUnique(*session.GetPath(), sessionID)
	} else {
		return errors.New("session not opened")
	}
}

func (r *blobRepo) Unlock(sessionID string) error {
	session, err := r.SessionManager().Get(sessionID)
	if err != nil {
		return err
	}
	return r.locks.Unlock(*session.GetPath())
}

func (r *blobRepo) SessionManager() server.SessionManager {
	return r.sessions
}

func (r *blobRepo) DataServerManager() server.DataServerManager {
	return r.servers
}

type repo struct {
	blobRepo repo2.BlobRepo
}

//var _ repo3.BlobRepo = (*repo)(nil)

var (
	r    repo
	once sync.Once
)

// Repo creates and returns the store client instance.
func Repo(servers server.DataServerManager, sessions server.SessionManager, locks server.LockManager) (repo2.Repo, error) {
	once.Do(func() {
		r = repo{
			blobRepo: newBlobRepo(servers, sessions, locks),
		}
	})

	return r, nil
}

func (r repo) BlobRepo() repo2.BlobRepo {
	return r.blobRepo
}

// Close closes the repo.
func (r repo) Close() error {
	return r.Close()
}
