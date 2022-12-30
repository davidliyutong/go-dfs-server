package server

import (
	"errors"
)

type lockManager struct {
	volume string
	locks  map[string]map[string]bool
}

type LockManager interface {
	Lock(path string, sessionID string) error
	LockUnique(path string, id string) error
	Unlock(path string) error
	GetLock(path string) ([]string, error)
}

func NewLockManager(volume string) LockManager {
	return &lockManager{
		volume: volume,
		locks:  make(map[string]map[string]bool),
	}
}

func (l *lockManager) Lock(path string, id string) error {
	locks, ok := l.locks[path]
	if ok {
		_, ok := locks[id]
		if ok {
			return errors.New("file already locked by this session")
		} else {
			locks[id] = true
			return nil
		}
	} else {
		locks = make(map[string]bool)
		locks[id] = true
		l.locks[path] = locks
		return nil
	}
}

func (l *lockManager) LockUnique(path string, id string) error {
	locks, err := l.GetLock(path)
	if err != nil {
		return err
	}
	if locks != nil {
		return errors.New("file already locked")
	} else {
		return l.Lock(path, id)
	}
}

func (l *lockManager) Unlock(path string) error {
	_, ok := l.locks[path]
	if ok {
		delete(l.locks, path)
		return nil
	} else {
		return errors.New("file not locked")
	}
}

func (l *lockManager) GetLock(path string) ([]string, error) {
	locks, ok := l.locks[path]
	if ok {
		keys := make([]string, 0, len(locks))
		for k := range locks {
			if locks[k] {
				keys = append(keys, k)
			}
		}
		return keys, nil
	} else {
		return nil, nil
	}

}
