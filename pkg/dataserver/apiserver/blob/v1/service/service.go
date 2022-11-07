package v1

import (
	"go-dfs-server/pkg/nameserver/apiserver/dataserver/v1/repo"
)

type Service interface {
}

type service struct {
	repo repo.Repo
}

var _ Service = (*service)(nil)

func NewService(repo repo.Repo) Service {
	return &service{repo}
}
