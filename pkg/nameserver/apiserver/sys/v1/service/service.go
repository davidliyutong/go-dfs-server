package v1

import (
	"go-dfs-server/pkg/nameserver/apiserver/sys/v1/repo"
	"go-dfs-server/pkg/nameserver/server"
)

type SysService interface {
	Info() (string, error)
	GetSession(id string) (server.Session, error)
	GetSessions() ([]string, error)
}

var _ SysService = (*sysService)(nil)

type sysService struct {
	repo repo.Repo
}

func (o *sysService) GetSessions() ([]string, error) {
	sessions := server.BlobSessionManager.ListSessions()
	res := make([]string, len(sessions))
	for idx, session := range sessions {
		res[idx] = *session.GetID()
	}
	return res, nil
}
func (o *sysService) GetSession(id string) (server.Session, error) {
	return server.BlobSessionManager.Get(id)
}

func (o *sysService) Info() (string, error) {
	return server.GlobalServerDesc.Opt.UUID, nil
}

func newSysService(repo repo.Repo) SysService {
	return &sysService{repo: repo}
}

type Service interface {
	NewSysService() SysService
}

type service struct {
	repo repo.Repo
}

var _ Service = (*service)(nil)

func NewService(repo repo.Repo) Service {
	return &service{repo: repo}
}

func (s *service) NewSysService() SysService {
	return newSysService(s.repo)
}
