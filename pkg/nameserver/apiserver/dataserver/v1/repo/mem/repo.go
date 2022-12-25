package mem

import (
	repoInterface "go-dfs-server/pkg/nameserver/apiserver/dataserver/v1/repo"
	"sync"
)

type repo struct {
	dataServerRepo repoInterface.DataServerRepo
}

var (
	r    repoInterface.Repo
	once sync.Once
)

func (r *repo) DataServerRepo() repoInterface.DataServerRepo {
	return r.dataServerRepo
}

func (r *repo) Close() error {
	return nil
}

func InitRepo() (repoInterface.Repo, error) {
	once.Do(func() {
		r = &repo{
			dataServerRepo: newDataServerRepo(),
		}
	})
	return r, nil
}
