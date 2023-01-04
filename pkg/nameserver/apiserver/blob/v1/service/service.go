package v1

import (
	"github.com/gin-gonic/gin"
	model "go-dfs-server/pkg/nameserver/apiserver/blob/v1/model"
	"go-dfs-server/pkg/nameserver/apiserver/blob/v1/repo"
	"mime/multipart"
)

type BlobService interface {
	Sync(sessionID string, blob model.BlobMetaData) (model.BlobMetaData, error)
	Ls(path string) (bool, []model.BlobMetaData, error)
	Mkdir(path string) error
	Open(path string, mode int) (model.BlobMetaData, error)
	Read(path string, chunkID int64, chunkOffset int64, size int64, c *gin.Context) (int64, error)
	Rm(path string, recursive bool) error
	Rmdir(path string) error
	Write(path string, chunkID int64, chunkOffset int64, size int64, version int64, file *multipart.FileHeader) ([]string, int64, error)
}

type blobService struct {
	repo repo.Repo
}

var _ BlobService = (*blobService)(nil)

func newBlobService(repo repo.Repo) BlobService {
	return &blobService{repo: repo}
}

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
