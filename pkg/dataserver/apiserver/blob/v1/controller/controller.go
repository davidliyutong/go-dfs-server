package v1

import (
	"github.com/gin-gonic/gin"
	"go-dfs-server/pkg/dataserver/apiserver/blob/v1/repo"
	srv "go-dfs-server/pkg/dataserver/apiserver/blob/v1/service"
)

type Controller interface {
	CreateChunk(c *gin.Context)
	CreateDirectory(c *gin.Context)
	CreateFile(c *gin.Context)
	DeleteChunk(c *gin.Context)
	DeleteDirectory(c *gin.Context)
	DeleteFile(c *gin.Context)
	LockFile(c *gin.Context)
	ReadChunk(c *gin.Context)
	ReadChunkMeta(c *gin.Context)
	ReadFileLock(c *gin.Context)
	ReadFileMeta(c *gin.Context)
	UnlockFile(c *gin.Context)
	WriteChunk(c *gin.Context)
}

type controller struct {
	srv srv.Service
}

func NewController(repo repo.Repo) Controller {
	return &controller{
		srv: srv.NewService(repo),
	}
}
