package v1

import (
	"github.com/gin-gonic/gin"
	model "go-dfs-server/pkg/nameserver/apiserver/blob/v1/model"
	"go-dfs-server/pkg/nameserver/apiserver/blob/v1/repo"
	v1 "go-dfs-server/pkg/nameserver/server"
	"mime/multipart"
)

type BlobService interface {
	Close(sessionID string) error
	Flush(sessionID string) error
	GetLock(sessionID string) ([]string, error)
	GetFileMeta(path string) (model.BlobMetaData, error)
	Lock(sessionID string) error
	Ls(path string) ([]model.BlobMetaData, error)
	Mkdir(path string) error
	Open(path string, mode int) (v1.Session, error)
	Read(sessionID string, size int64, c *gin.Context) (int64, error)
	Rm(path string, recursive bool) error
	Rmdir(path string) error
	Seek(sessionID string, offset int64, whence int) (int64, error)
	Truncate(sessionID string, size int64) error
	Unlock(sessionID string) error
	Write(sessionID string, syncWrite bool, file *multipart.FileHeader) (int64, error)
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
