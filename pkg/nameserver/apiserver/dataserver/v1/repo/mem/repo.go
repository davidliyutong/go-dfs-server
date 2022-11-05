package mem

import (
	repoInterface "go-dfs-server/pkg/nameserver/apiserver/dataserver/v1/repo"
	"sync"
)

type repo struct {
	dataserverRepo repoInterface.DataserverRepo
}

var (
	r    repoInterface.Repo
	once sync.Once
)

func InitRepo() (repoInterface.Repo, error) {
	once.Do(func() {
		r = repo{
			dataserverRepo: newDataserverRepo(),
		}
	})
	return r, nil
}

func (r repo) DataserverRepo() repoInterface.DataserverRepo {
	return r.dataserverRepo
}

func (r repo) Close() error {
	return nil
}
