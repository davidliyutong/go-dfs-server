package v1

import (
	"go-dfs-server/pkg/nameserver/apiserver/blob/v1/repo"
	srv "go-dfs-server/pkg/nameserver/apiserver/blob/v1/service"
)

type Controller interface {
}

type controller struct {
	srv srv.Service
}

func NewController(repo repo.Repo) Controller {
	return &controller{
		srv: srv.NewService(repo),
	}
}
