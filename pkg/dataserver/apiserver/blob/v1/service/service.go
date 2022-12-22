package v1

import (
	"go-dfs-server/pkg/dataserver/apiserver/blob/v1/repo"
)

type Service interface {
	NewBlobService() BlobService
}

type service struct {
	repo repo.Repo
}

var _ Service = (*service)(nil)

func NewService(repo repo.Repo) Service {
	return &service{repo}
}

func (s *service) NewBlobService() BlobService {
	return newBlobService(s.repo)
}
