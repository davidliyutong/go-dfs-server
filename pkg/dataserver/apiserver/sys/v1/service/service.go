package v1

import (
	"go-dfs-server/pkg/nameserver/apiserver/sys/v1/repo"
)

type SysService interface {
	//Info() (string, error)
}

var _ SysService = (*sysService)(nil)

type sysService struct {
	repo repo.Repo
}

func newSysService(repo repo.Repo) SysService {
	return &sysService{repo: repo}
}

type Service interface {
	NewSysService(repo repo.Repo) SysService
}

type service struct {
	repo repo.Repo
}

var _ Service = (*service)(nil)

func NewService(repo repo.Repo) Service {
	return &service{repo}
}

func (*service) NewSysService(repo repo.Repo) SysService {
	return newSysService(repo)
}
